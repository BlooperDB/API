package nodes

import (
	"errors"
	"strconv"

	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/storage"
	"github.com/BlooperDB/API/utils"
	"github.com/graphql-go/graphql"
	"github.com/wuman/firebase-server-sdk-go"
)

var enumVote = graphql.NewEnum(
	graphql.EnumConfig{
		Name: "Vote",
		Values: graphql.EnumValueConfigMap{
			"UP":   &graphql.EnumValueConfig{},
			"DOWN": &graphql.EnumValueConfig{},
			"NONE": &graphql.EnumValueConfig{},
		},
	},
)

var graphTag = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tag",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)

var graphComment = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"user": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"revisionId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"message": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)

var graphRevision = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Revision",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"revision": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"changes": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"blueprintId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"blueprint": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"thumbsUp": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"thumbsDown": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"userVote": &graphql.Field{
				Type: graphql.NewNonNull(enumVote),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					switch utils.Source(p, "userVote").(int) {
					default:
						return "NONE", nil
					case 1:
						return "UP", nil
					case 2:
						return "DOWN", nil
					}
				},
			},
			"comments": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphComment)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToComments(utils.Source(p, "_db").(*db.Revision).GetComments()), nil
				},
			},
			"version": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"thumbnail": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"render": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)

var graphBlueprint = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Blueprint",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"user": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"latestRevision": &graphql.Field{
				Type: graphql.NewNonNull(graphRevision),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToRevision(utils.Source(p, "_db").(*db.Blueprint).GetLatestRevision(), db.GetAuthUserGraphQL(p)), nil
				},
			},
			"revisions": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphRevision)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToRevisions(utils.Source(p, "_db").(*db.Blueprint).GetRevisions(), db.GetAuthUserGraphQL(p)), nil
				},
			},
			"revision": &graphql.Field{
				Type: graphql.NewNonNull(graphRevision),
				Args: graphql.FieldConfigArgument{
					"revision": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToRevision(utils.Source(p, "_db").(*db.Blueprint).GetRevision(uint(p.Args["revision"].(int))), db.GetAuthUserGraphQL(p)), nil
				},
			},
			"tags": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphTag)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToTags(db.GetTagsFromBlueprint(utils.Source(p, "id").(uint))), nil
				},
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"thumbnail": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)

var interfaceUserData = graphql.NewInterface(
	graphql.InterfaceConfig{
		Name: "UserData",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"avatar": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"blueprints": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphBlueprint)),
			},
		},
	},
)

var graphPublicUser = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PublicUser",
		Interfaces: []*graphql.Interface{
			interfaceUserData,
		},
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"avatar": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"blueprints": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphBlueprint)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToBlueprints(utils.Source(p, "_db").(*db.User).GetUserBlueprints()), nil
				},
			},
		},
	},
)

var graphPrivateUser = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PrivateUser",
		Interfaces: []*graphql.Interface{
			interfaceUserData,
		},
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"avatar": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"blueprints": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphBlueprint)),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToBlueprints(utils.Source(p, "_db").(*db.User).GetUserBlueprints()), nil
				},
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)

var unionUser = graphql.NewUnion(
	graphql.UnionConfig{
		Name: "User",
		Types: []*graphql.Object{
			graphPublicUser,
			graphPrivateUser,
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			_, ok := p.Value.(map[string]interface{})["email"]

			if ok {
				return graphPrivateUser
			}

			return graphPublicUser
		},
	},
)

var enumBlueprintOrder = graphql.NewEnum(
	graphql.EnumConfig{
		Name: "BlueprintOrder",
		Values: graphql.EnumValueConfigMap{
			"NORMAL":  &graphql.EnumValueConfig{},
			"POPULAR": &graphql.EnumValueConfig{},
			"TOP":     &graphql.EnumValueConfig{},
			"NEW":     &graphql.EnumValueConfig{},
		},
	},
)

var graphQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"blueprint": &graphql.Field{
				Type:        graphBlueprint,
				Description: "Retrieve blueprint by id.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToBlueprint(db.GetBlueprintById(uint(p.Args["id"].(int)))), nil
				},
			},
			"blueprints": &graphql.Field{
				Type:        graphql.NewList(graphBlueprint),
				Description: "Retrieve blueprints.",
				Args: graphql.FieldConfigArgument{
					"order": &graphql.ArgumentConfig{
						Type:         enumBlueprintOrder,
						DefaultValue: "NORMAL",
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"count": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 20,
					},
					"search": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"ascending": &graphql.ArgumentConfig{
						Type:         graphql.Boolean,
						DefaultValue: false,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					count := utils.MinMax(1, p.Args["count"].(int), 100)
					return dbToBlueprints(db.FindBlueprintsDynamic(p.Args["search"].(string), p.Args["offset"].(int), count, p.Args["order"].(string), p.Args["ascending"].(bool))), nil
				},
			},
			"tags": &graphql.Field{
				Type:        graphql.NewList(graphTag),
				Description: "Retrieve tags. By default will return popular tags.",
				Args: graphql.FieldConfigArgument{
					"autocomplete": &graphql.ArgumentConfig{
						Type:         graphql.String,
						Description:  "Autocomplete tags",
						DefaultValue: "",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if p.Args["autocomplete"].(string) != "" {
						return dbToTags(db.AutocompleteTag(p.Args["autocomplete"].(string))), nil
					}

					return dbToTags(db.PopularTags()), nil
				},
			},
			"revision": &graphql.Field{
				Type:        graphRevision,
				Description: "Retrieve revision by id.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToRevision(db.GetRevisionById(uint(p.Args["id"].(int))), db.GetAuthUserGraphQL(p)), nil
				},
			},
			"user": &graphql.Field{
				Type:        unionUser,
				Description: "Retrieve self or user by id.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"self": &graphql.ArgumentConfig{
						Type:         graphql.Boolean,
						DefaultValue: false,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if p.Args["self"].(bool) {
						user := db.GetAuthUserGraphQL(p)
						if user != nil {
							return dbToPrivateUser(user), nil
						} else {
							return nil, errors.New("invalid token")
						}
					}

					return dbToPublicUser(db.GetUserById(uint(p.Args["id"].(int)))), nil
				},
			},
			"comment": &graphql.Field{
				Type:        graphComment,
				Description: "Retrieve comment by id.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return dbToComment(db.GetCommentById(uint(p.Args["id"].(int)))), nil
				},
			},
		},
	},
)

var graphSignInResponse = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "SignIn",
		Fields: graphql.Fields{
			"blooperToken": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"firstLogin": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
	},
)

var graphMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"signIn": &graphql.Field{
				Type:        graphSignInResponse,
				Description: "Log in as a user.",
				Args: graphql.FieldConfigArgument{
					"firebaseToken": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					auth, _ := firebase.GetAuth()
					decodedToken, err := auth.VerifyIDToken(p.Args["firebaseToken"].(string))

					if err != nil {
						return nil, errors.New("user token invalid: " + err.Error())
					}

					_, found := decodedToken.UID()

					if !found {
						return nil, errors.New("user token invalid")
					}

					user, firstLogin := db.SignIn(decodedToken)

					return map[string]interface{}{
						"blooperToken": user.BlooperToken,
						"firstLogin":   firstLogin,
					}, nil
				},
			},
			"rateRevision": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Rate a revision.",
				Args: graphql.FieldConfigArgument{
					"revision": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"vote": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(enumVote),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					revision := db.GetRevisionById(uint(p.Args["revision"].(int)))

					if revision == nil {
						return nil, errors.New("revision not found")
					}

					rating := db.FindRating(user.ID, revision.ID)

					vote := p.Args["vote"].(string)

					if vote == "NONE" {
						if rating.ID == 0 {
							return nil, errors.New("rating not found")
						}

						rating.Delete()
					} else {
						thumbsUp := true

						if vote == "DOWN" {
							thumbsUp = false
						}

						rating.UserID = user.ID
						rating.RevisionID = revision.ID
						rating.ThumbsUp = thumbsUp
						rating.DeletedAt = nil
						rating.Save()
					}

					return true, nil
				},
			},
			"addRevision": &graphql.Field{
				Type:        graphRevision,
				Description: "Add a revision.",
				Args: graphql.FieldConfigArgument{
					"blueprintId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"changes": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"blueprint": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					blueprint := db.GetBlueprintById(uint(p.Args["blueprintId"].(int)))

					if blueprint == nil {
						return nil, errors.New("blueprint not found")
					}

					if blueprint.UserID != user.ID {
						return nil, errors.New("unable to mutate this blueprint")
					}

					i := blueprint.IncrementAndGetRevision()

					blueprintString := p.Args["blueprint"].(string)
					changes := p.Args["changes"].(string)

					bpVersion, _ := strconv.Atoi(blueprintString[0:1])

					sha265 := utils.SHA265(blueprintString)

					if db.FindRevisionByChecksum(sha265) != nil {
						return nil, errors.New("blueprint already exists")
					}

					revision := &db.Revision{
						BlueprintID:       blueprint.ID,
						Revision:          i,
						Changes:           changes,
						BlueprintVersion:  bpVersion,
						BlueprintChecksum: sha265,
					}

					revision.Save()

					storage.SaveRevision(revision.ID, blueprintString)
					go storage.RenderAndSaveAndUpdateBlueprint(blueprintString, revision)

					return dbToRevision(revision, db.GetAuthUserGraphQL(p)), nil
				},
			},
			"updateRevision": &graphql.Field{
				Type:        graphRevision,
				Description: "Update a revision.",
				Args: graphql.FieldConfigArgument{
					"revision": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"changes": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					revision := db.GetRevisionById(uint(p.Args["revision"].(int)))

					if revision == nil {
						return nil, errors.New("revision not found")
					}

					blueprint := revision.GetBlueprint()

					if blueprint.UserID != user.ID {
						return nil, errors.New("unable to mutate this blueprint")
					}

					revision.Changes = p.Args["changes"].(string)
					revision.Save()

					return dbToRevision(revision, db.GetAuthUserGraphQL(p)), nil
				},
			},
			"deleteRevision": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a revision.",
				Args: graphql.FieldConfigArgument{
					"revision": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					revision := db.GetRevisionById(uint(p.Args["revision"].(int)))

					if revision == nil {
						return nil, errors.New("revision not found")
					}

					blueprint := revision.GetBlueprint()

					if blueprint.UserID != user.ID {
						return nil, errors.New("unable to mutate this blueprint")
					}

					revision.Delete()

					if blueprint.CountRevisions() == 0 {
						blueprint.Delete()
					}

					return true, nil
				},
			},
			"addBlueprint": &graphql.Field{
				Type:        graphBlueprint,
				Description: "Add a blueprint.",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"blueprint": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					name := p.Args["name"].(string)
					description := p.Args["description"].(string)
					blueprintString := p.Args["blueprint"].(string)
					tags := p.Args["tags"].([]string)

					sha265 := utils.SHA265(blueprintString)

					if db.FindRevisionByChecksum(sha265) != nil {
						return nil, errors.New("blueprint already exists")
					}

					blueprint := &db.Blueprint{
						UserID:       user.ID,
						Name:         name,
						Description:  description,
						LastRevision: 1,
					}

					blueprint.Save()

					bpVersion, _ := strconv.Atoi(blueprintString[0:1])

					revision := &db.Revision{
						BlueprintID:       blueprint.ID,
						Revision:          blueprint.LastRevision,
						Changes:           "",
						BlueprintVersion:  bpVersion,
						BlueprintChecksum: sha265,
					}

					revision.Save()

					storage.SaveRevision(revision.ID, blueprintString)
					go storage.RenderAndSaveAndUpdateBlueprint(blueprintString, revision)

					for _, tag := range tags {

						t := db.GetTagByName(tag)

						if t == nil {
							t = &db.Tag{
								Name: tag,
							}

							t.Save()
						}

						bt := db.BlueprintTag{
							BlueprintId: blueprint.ID,
							TagId:       t.ID,
						}

						bt.Save()
					}

					return dbToBlueprint(blueprint), nil
				},
			},
			"updateBlueprint": &graphql.Field{
				Type:        graphBlueprint,
				Description: "Update a blueprint.",
				Args: graphql.FieldConfigArgument{
					"blueprintId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					name := p.Args["name"].(string)
					description := p.Args["description"].(string)
					tags := p.Args["blueprint"].([]string)

					blueprint := db.GetBlueprintById(uint(p.Args["blueprintId"].(int)))

					if blueprint == nil {
						return nil, errors.New("blueprint not found")
					}

					if blueprint.UserID != user.ID {
						return nil, errors.New("unable to mutate this blueprint")
					}

					for _, t := range blueprint.GetTags() {
						bt := db.BlueprintTag{
							BlueprintId: blueprint.ID,
							TagId:       t.ID,
						}

						bt.Delete()
					}

					for _, tag := range tags {
						t := &db.Tag{
							Name: tag,
						}

						t.Save()

						bt := db.BlueprintTag{
							BlueprintId: blueprint.ID,
							TagId:       t.ID,
						}

						bt.Save()
					}

					blueprint.Name = name
					blueprint.Description = description

					blueprint.Save()

					return dbToBlueprint(blueprint), nil
				},
			},
			"deleteBlueprint": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a blueprint.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					blueprint := db.GetBlueprintById(uint(p.Args["id"].(int)))

					if blueprint == nil {
						return nil, errors.New("blueprint not found")
					}

					if blueprint.UserID != user.ID {
						return nil, errors.New("unable to mutate this blueprint")
					}

					blueprint.Delete()

					return true, nil
				},
			},
			"addComment": &graphql.Field{
				Type:        graphComment,
				Description: "Add a comment.",
				Args: graphql.FieldConfigArgument{
					"revision": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"message": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					comment := &db.Comment{
						RevisionID: uint(p.Args["revision"].(int)),
						UserID:     user.ID,
						Message:    p.Args["message"].(string),
					}

					comment.Save()

					return dbToComment(comment), nil
				},
			},
			"updateComment": &graphql.Field{
				Type:        graphComment,
				Description: "Update a comment.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"message": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					comment := db.GetCommentById(uint(p.Args["id"].(int)))

					if comment == nil {
						return nil, errors.New("comment not found")
					}

					if comment.ID != user.ID {
						return nil, errors.New("unable to mutate this comment")
					}

					comment.Message = p.Args["message"].(string)
					comment.Save()

					return dbToComment(comment), nil
				},
			},
			"deleteComment": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a comment.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := db.GetAuthUserGraphQL(p)

					if user == nil {
						return nil, errors.New("invalid token")
					}

					comment := db.GetCommentById(uint(p.Args["id"].(int)))

					if comment == nil {
						return nil, errors.New("comment not found")
					}

					if comment.ID != user.ID {
						return nil, errors.New("unable to mutate this comment")
					}

					comment.Delete()

					return true, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    graphQuery,
		Mutation: graphMutation,
	},
)

func InitializeGraphs() {
}

func GetSchema() *graphql.Schema {
	return &schema
}

func dbToBlueprint(blueprint *db.Blueprint) interface{} {
	if blueprint == nil {
		return nil
	}

	return map[string]interface{}{
		"_db":         blueprint,
		"id":          blueprint.ID,
		"user":        blueprint.UserID,
		"name":        blueprint.Name,
		"description": blueprint.Description,
		"createdAt":   blueprint.CreatedAt,
		"updatedAt":   blueprint.UpdatedAt,
		"thumbnail":   storage.PublicURL + "/" + storage.BlueprintRenderBucket + "/" + blueprint.GetThumbnail() + "-thumbnail.png",
	}
}

func dbToBlueprints(blueprints []*db.Blueprint) []interface{} {
	var result []interface{}

	for _, blueprint := range blueprints {
		result = append(result, dbToBlueprint(blueprint))
	}

	return result
}

func dbToTag(tag *db.Tag) interface{} {
	if tag == nil {
		return nil
	}

	return map[string]interface{}{
		"_db":  tag,
		"name": tag.Name,
	}
}

func dbToTags(tags []*db.Tag) []interface{} {
	var result []interface{}

	for _, tag := range tags {
		result = append(result, dbToTag(tag))
	}

	return result
}

func dbToRevision(revision *db.Revision, user *db.User) interface{} {
	if revision == nil {
		return nil
	}

	if revision.DeletedAt != nil {
		return nil
	}

	ratings := revision.GetRatings()
	thumbsUp, thumbsDown, userVote := 0, 0, 0

	for _, rating := range ratings {
		if rating.ThumbsUp {
			thumbsUp++
		} else {
			thumbsDown++
		}

		if user != nil && user.ID == rating.UserID {
			if rating.ThumbsUp {
				userVote = 1
			} else {
				userVote = 2
			}
		}
	}

	baseRenderStorageURL := storage.PublicURL + "/" + storage.BlueprintRenderBucket + "/" + revision.BlueprintChecksum

	return map[string]interface{}{
		"_db":         revision,
		"id":          revision.ID,
		"revision":    revision.Revision,
		"changes":     revision.Changes,
		"createdAt":   revision.CreatedAt,
		"updatedAt":   revision.UpdatedAt,
		"blueprintId": revision.BlueprintID,
		"blueprint":   storage.PublicURL + "/" + storage.BlueprintStringBucket + "/" + storage.RevisionToString(revision.ID),
		"thumbsUp":    thumbsUp,
		"thumbsDown":  thumbsDown,
		"userVote":    userVote,
		"version":     revision.BlueprintVersion,
		"thumbnail":   baseRenderStorageURL + "-thumbnail.png",
		"render":      baseRenderStorageURL + ".png",
	}
}

func dbToRevisions(revisions []*db.Revision, user *db.User) []interface{} {
	var result []interface{}

	for _, revision := range revisions {
		result = append(result, dbToRevision(revision, user))
	}

	return result
}

func dbToComment(comment *db.Comment) interface{} {
	if comment == nil {
		return nil
	}

	return map[string]interface{}{
		"_db":        comment,
		"id":         comment.ID,
		"user":       comment.UserID,
		"createdAt":  comment.CreatedAt,
		"updatedAt":  comment.UpdatedAt,
		"message":    comment.Message,
		"revisionId": comment.RevisionID,
	}
}

func dbToComments(comments []*db.Comment) []interface{} {
	var result []interface{}

	for _, comment := range comments {
		result = append(result, dbToComment(comment))
	}

	return result
}

func dbToPublicUser(user *db.User) interface{} {
	if user == nil {
		return nil
	}

	return map[string]interface{}{
		"_db":      user,
		"id":       user.ID,
		"username": user.Username,
		"avatar":   user.Avatar,
	}
}

func dbToPrivateUser(user *db.User) interface{} {
	u := dbToPublicUser(user)

	if u == nil {
		return nil
	}

	u.(map[string]interface{})["email"] = user.Email
	u.(map[string]interface{})["createdAt"] = user.CreatedAt
	u.(map[string]interface{})["updatedAt"] = user.UpdatedAt

	return u
}
