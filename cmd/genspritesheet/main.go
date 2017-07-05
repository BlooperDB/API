package main

import (
	"fmt"

	"strings"

	"strconv"

	"encoding/json"

	"github.com/yuin/gopher-lua"
)

type EntityData struct {
	Type    string
	Name    string
	Icon    string
	Picture Picture
}

type Picture struct {
	States map[string][]Layer
}

type Layer struct {
	File   string
	Width  int
	Height int
	Shift  Shift
}

type Shift struct {
	X float64
	Y float64
}

func main() {

	fmt.Println("Generating spritesheet...")

	L := lua.NewState(lua.Options{
		IncludeGoStackTrace: true,
	})

	L.OpenLibs()

	L.SetGlobal("processExtension", L.NewFunction(processExtension))

	L.DoString(`
		data = {
			is_demo=false,
			raw={
				player={
					player={
						animations={}
					}
				},
				explosion={
					explosion={animations={{}}},
					["big-explosion"]={animations={{}}}
				}
			}
		}

		function data.extend(self, extra)
			processExtension(extra)
		end

		defines = {direction={north=0,east=2,south=4,west=6}}
	`)

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-rail-pictures.lua")
	L.PreloadModule("prototypes.entity.demo-rail-pictures", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-pipecovers.lua")
	L.PreloadModule("prototypes.entity.demo-pipecovers", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-transport-belt-pictures.lua")
	L.PreloadModule("prototypes.entity.demo-transport-belt-pictures", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-circuit-connector-sprites.lua")
	L.PreloadModule("prototypes.entity.demo-circuit-connector-sprites", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-player-animations.lua")
	L.PreloadModule("prototypes.entity.demo-player-animations", func(state *lua.LState) int { return 0 })

	L.PreloadModule("util", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/transport-belt-pictures.lua")
	L.PreloadModule("prototypes.entity.transport-belt-pictures", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/assemblerpipes.lua")
	L.PreloadModule("prototypes.entity.assemblerpipes", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/laser-sounds.lua")
	L.PreloadModule("prototypes.entity.laser-sounds", func(state *lua.LState) int { return 0 })

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-gunshot-sounds.lua")
	L.PreloadModule("prototypes.entity.demo-gunshot-sounds", func(state *lua.LState) int { return 0 })

	L.PreloadModule("prototypes.entity.pump-connector", func(state *lua.LState) int { return 0 })

	playerAnimations := L.NewTable()
	playerAnimations.RawSetString("level1", L.NewTable())
	playerAnimations.RawSetString("level2addon", L.NewTable())
	playerAnimations.RawSetString("level3addon", L.NewTable())
	L.SetGlobal("playeranimations", playerAnimations)

	L.SetGlobal("module", L.NewFunction(func(state *lua.LState) int {
		return 0
	}))

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/core/lualib/util.lua")

	L.DoString(`
		_G.util = {
			distance=_G.distance,
			findfirstentity=_G.findfirstentity,
			positiontostr=_G.positiontostr,
			formattime=_G.formattime,
			moveposition=_G.moveposition,
			oppositedirection=_G.oppositedirection,
			ismoduleavailable=_G.ismoduleavailable,
			multiplystripes=_G.multiplystripes,
			by_pixel=_G.by_pixel,
			format_number=_G.format_number,
			increment=_G.increment,
			table={deepcopy=_G.table.deepcopy}
		}

		function util.table.deepcopy(extra)
			return extra
		end
	`)

	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/demo-entities.lua")
	executeFile(L, "src/github.com/BlooperDB/API/factorio/data/base/prototypes/entity/entities.lua")

	defer L.Close()
}

func dumpTable(table *lua.LTable, indent int) {
	table.ForEach(func(key lua.LValue, val lua.LValue) {
		if !strings.HasPrefix(key.String(), "_") {
			if val.Type() == lua.LTTable {
				fmt.Println(spaces(indent) + key.String() + ": ")
				dumpTable(val.(*lua.LTable), indent+2)
			} else {
				fmt.Println(spaces(indent) + key.String() + ": " + val.String())
			}
		}
	})
}

func spaces(count int) string {
	spaces := ""
	for i := 0; i < count; i++ {
		spaces = spaces + " "
	}
	return spaces
}

func executeFile(L *lua.LState, file string) {
	if err := L.DoFile(file); err != nil {
		panic(err)
	}
}

func processExtension(state *lua.LState) int {
	tbl := state.ToTable(1)

	result := make(map[string]EntityData)

	tbl.ForEach(func(key lua.LValue, val lua.LValue) {
		entity := val.(*lua.LTable)

		entityIcon := entity.RawGetString("icon").String()

		if entityIcon == "nil" {
			return
		}

		entityType := entity.RawGetString("type").String()

		if entityType == "player" ||
			entityType == "corpse" ||
			entityType == "character-corpse" ||
			entityType == "fish" ||
			entityType == "combat-robot" ||
			entityType == "logistic-robot" ||
			entityType == "construction-robot" {
			return
		}

		entityName := entity.RawGetString("name").String()

		picture := Picture{
			States: make(map[string][]Layer),
		}

		animationData := entity.RawGetString("animation")
		animationsData := entity.RawGetString("animations")

		if animationData.String() != "nil" {
			layers := animationData.(*lua.LTable)
			processPictures(layers, &picture, "default")
		} else if animationsData.String() != "nil" {
			processPictures(animationsData.(*lua.LTable), &picture, "default")
		}

		pictureData := entity.RawGetString("picture")

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("pictures")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("picture_on")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("structure")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("picture_safe")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("base")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("sprite")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("sprites")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("power_on_animation")
		}

		if pictureData.String() == "nil" {
			pictureData = entity.RawGetString("off_animation")
		}

		if pictureData.String() != "nil" {
			processPictures(pictureData.(*lua.LTable), &picture, "default")
		}

		if pictureData.String() == "nil" {
			if entity.RawGetString("horizontal_animation").String() != "nil" {
				processPictures(entity.RawGetString("horizontal_animation").(*lua.LTable), &picture, "horizontal")
				processPictures(entity.RawGetString("vertical_animation").(*lua.LTable), &picture, "vertical")
			}
		}

		if len(picture.States) == 0 {
			fmt.Println(entityType + ": " + entityName)
		}

		result[entityName] = EntityData{
			Type:    entityType,
			Name:    entityName,
			Icon:    entityIcon,
			Picture: picture,
		}
	})

	j, _ := json.Marshal(result)
	fmt.Println(string(j))

	return 0
}

func processPictures(pTable *lua.LTable, picture *Picture, defaultName string) {
	if pTable.RawGetString("filename").String() != "nil" {
		// There is only 1 state with 1 layer
		picture.States[defaultName] = make([]Layer, 1)
		picture.States[defaultName][0] = processLayer(pTable)
	} else {
		potentialLayers := pTable.RawGetString("layers")
		if potentialLayers.String() != "nil" {
			// There is only 1 state with multiple layers
			picture.States[defaultName] = processLayers(potentialLayers.(*lua.LTable))
		} else {
			// There are multiple states
			pTable.ForEach(func(state lua.LValue, stateData lua.LValue) {
				stateTable := stateData.(*lua.LTable)
				if stateTable.RawGetString("filename").String() != "nil" {
					// There are multiple states with 1 layer
					picture.States[state.String()] = make([]Layer, 1)
					picture.States[state.String()][0] = processLayer(stateTable)
				} else {
					potentialLayers := stateTable.RawGetString("layers")
					if potentialLayers.String() != "nil" {
						// There are multiple states with multiple layers
						picture.States[state.String()] = processLayers(potentialLayers.(*lua.LTable))
					} else {
						potentialSheet := stateTable.RawGetString("sheet")
						if potentialSheet.String() != "nil" {
							// There are multiple states with 1 layer
							picture.States[state.String()] = make([]Layer, 1)
							picture.States[state.String()][0] = processLayer(potentialSheet.(*lua.LTable))
						} else {
							// There are multiple states with multiple layers with multiple options (we pick first)
							first := stateTable.RawGetInt(1)

							if first.String() == "nil" {
								return
							}

							picture.States[state.String()] = processLayers(first.(*lua.LTable).RawGetString("layers").(*lua.LTable))
						}
					}
				}
			})

		}

	}
}

func processLayers(pTable *lua.LTable) []Layer {
	layers := make([]Layer, pTable.Len())

	pTable.ForEach(func(pos lua.LValue, v lua.LValue) {
		i, _ := strconv.Atoi(pos.String())
		layers[i-1] = processLayer(v.(*lua.LTable))
	})

	return layers
}

func processLayer(pTable *lua.LTable) Layer {
	width, _ := strconv.Atoi(pTable.RawGetString("width").String())
	height, _ := strconv.Atoi(pTable.RawGetString("height").String())

	ShiftX, ShiftY := float64(0), float64(0)

	if pTable.RawGetString("shift").String() != "nil" {
		ShiftX, _ = strconv.ParseFloat(pTable.RawGetString("shift").(*lua.LTable).RawGetInt(1).String(), 32)
		ShiftY, _ = strconv.ParseFloat(pTable.RawGetString("shift").(*lua.LTable).RawGetInt(2).String(), 32)
	}

	return Layer{
		File:   pTable.RawGetString("filename").String(),
		Width:  width,
		Height: height,
		Shift: Shift{
			X: ShiftX,
			Y: ShiftY,
		},
	}
}
