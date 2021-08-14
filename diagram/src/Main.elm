module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Debug
import String
import List
import List.Extra
import Array

-- MAIN

main =
  Browser.element
      { init = init
      , view = view
      , update = update
      , subscriptions = subscriptions
      }

-- MODEL

type alias Caption =
    { text: String
    , size: String
    , position: String
    }

toRow: Int -> Caption -> Html Msg
toRow id cap =
    tr []
        [ td [] [text cap.text]
        , td [] [text cap.size]
        , td [] [text cap.position]
        , td [] [button [onClick (EditCap id)] [ text "edit" ]]
        , td [] [button [onClick (CopyCap id)] [ text "copy" ]]
        , td [] [button [onClick (RemoveCap id)] [ text "del" ]]
        ]

type alias Model =
    { image : String
    , gridx : String
    , gridy : String
    , oldGridx : String
    , oldGridy : String
    , text : String
    , size : String
    , position: String
    , oldText : String
    , oldSize : String
    , oldPosition: String
    , captions: List Caption
    }

init : String -> (Model, Cmd Msg)
init image =
    ( { image = image
      , gridx = ""
      , gridy = ""
      , oldGridx = ""
      , oldGridy = ""
      , text = ""
      , size = ""
      , position = ""
      , oldText = ""
      , oldSize = ""
      , oldPosition = ""
      , captions = []
      }
    , Cmd.none
    )

-- UPDATE

type Msg
    = Grid
    | ClearGrid
    | AddCap
    | EditCap Int
    | CopyCap Int
    | RemoveCap Int
    | GotImg (Result Http.Error String)
    | ChangeX String
    | ChangeY String
    | ChangeText String
    | ChangeSize String
    | ChangePosition String

popCap : Model -> Int -> Maybe Model
popCap model id =
    let myCap = Array.get id (Array.fromList model.captions)
    in case myCap of
            Nothing -> Nothing
            Just cap ->
                    Just {model
                        | captions = List.Extra.removeAt id model.captions
                        , text = cap.text
                        , size = cap.size
                        , position = cap.position
                    }

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        CopyCap id ->
            let myCap = Array.get id (Array.fromList model.captions)
            in case myCap of
                    Nothing ->
                        (model, Cmd.none)
                    Just cap ->
                        let mod = {model
                                      | text = cap.text
                                      , size = cap.size
                                      , position = cap.position
                                  }
                        in (mod, addCaption mod)
        EditCap id ->
            let newMod = popCap model id
            in case newMod of
                   Nothing -> (model, Cmd.none)
                   Just mod -> ( mod, addCaption mod )
        RemoveCap id ->
            let newMod = 
                    {model | captions =
                        List.Extra.removeAt id model.captions
                    }
            in ( newMod ,
                  addCaption newMod )
        AddCap ->
            if model.text == "" ||
                model.size == "" ||
                    model.position == ""
            then
                (model, Cmd.none )
            else
                let newMod =
                        { model | captions =
                              { text = model.text
                              , size = model.size
                              , position = model.position
                              } :: model.captions
                        , oldText = model.text
                        , oldSize = model.size
                        , oldPosition = model.position
                        , text = ""
                        , size = ""
                        , position = ""
                        }
                    in ( newMod , addCaption newMod )
        ChangeText newText ->
            ( { model | text = newText }, Cmd.none )
        ChangeSize newText ->
            ( { model | size = newText }, Cmd.none )
        ChangePosition newText ->
            ( { model | position = newText }, Cmd.none )
        ChangeX newX ->
            ( { model | gridx = newX }, Cmd.none )
        ChangeY newY ->
            ( { model | gridy = newY }, Cmd.none )
        Grid ->
            if model.gridx == "" ||
                model.gridy == ""
            then
                (model, Cmd.none)
            else
                let newMod = 
                        { model | image = model.image
                        , oldGridx = model.gridx
                        , oldGridy = model.gridy
                        , gridx = ""
                        , gridy = ""
                        }
                    in ( newMod, addGrid newMod)
        ClearGrid ->
                let newMod = 
                        { model | image = model.image
                        , oldGridx = ""
                        , oldGridy = ""
                        , gridx = ""
                        , gridy = ""
                        }
                    in ( newMod, addGrid newMod)
        GotImg result ->
            case result of
                Ok img ->
                    ( {model | image = img}, Cmd.none)
                Err _ ->
                    (model, Cmd.none)

-- VIEW

size : Int -- size in px of the input boxes
size = 50

view : Model -> Html Msg
view model =
  div []
    [ img [src model.image] []
    , div []
        [ input [ placeholder "grid h"
                , value model.gridx
                , style "width" (String.fromInt (2*size) ++ "px"), onInput ChangeX ] []
        , input [ placeholder "grid v"
                , value model.gridy
                , style "width" (String.fromInt (2*size) ++ "px"), onInput ChangeY ] []
        , button [ onClick Grid ] [ text "grid" ]
        , button [ onClick ClearGrid ] [ text "clear" ]
        ]
    , div []
        [ input [ placeholder "lx", style "width" (String.fromInt size ++ "px") ] []
        , input [ placeholder "uy", style "width" (String.fromInt size ++ "px") ] []
        , input [ placeholder "rx", style "width" (String.fromInt size ++ "px") ] []
        , input [ placeholder "by", style "width" (String.fromInt size ++ "px") ] []
        , button [] [ text "crop" ]
        ]
    , div []
        [ input [ placeholder "Text"
                , value model.text
                , style "width" (String.fromInt (4*size//3) ++ "px")
                , onInput ChangeText ] []
        , input [ placeholder "Size"
                , value model.size
                , style "width" (String.fromInt (4*size//3) ++ "px")
                , onInput ChangeSize ] []
        , input [ placeholder "Position"
                , value model.position
                , style "width" (String.fromInt (4*size//3 + 1) ++ "px")
                , onInput ChangePosition ] []
        , button [onClick AddCap] [ text "add caption" ]
        ]
    , table []
        ([ thead []
               [ th [] [text "Text"]
               , th [] [text "Size"]
               , th [] [text "x,y"]
               ]
         ]
             ++ List.indexedMap toRow model.captions
        )
    ]
      
-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none

-- HTTP
gridStr : Model -> String
gridStr model =
        "grid=" ++ model.oldGridx ++ "," ++ model.oldGridy

capStr : Model -> String
capStr model =
    List.foldr
    (\cap str -> str ++ cap.text ++ "," ++ cap.size ++ "," ++ cap.position ++ "&cap=")
    "cap=" model.captions

addGrid : Model -> Cmd Msg
addGrid model =
    Http.get
        { url = "http://localhost:8080/req?" ++ (gridStr model) ++ "&" ++ (capStr model)
        , expect = Http.expectString GotImg
        }

addCaption : Model -> Cmd Msg
addCaption model =
    Http.get
        { url = "http://localhost:8080/req?" ++ (gridStr model) ++ "&" ++ (capStr model)
        , expect = Http.expectString GotImg
        }
