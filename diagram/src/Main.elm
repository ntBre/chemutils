module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Debug
import String
import List

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
    {text: String
    ,size: String
    ,position: String
    }

toRow: Caption -> Html Msg
toRow cap =
    tr []
        [ td [] [text cap.text]
        , td [] [text  cap.size]
        , td [] [text cap.position]
        ]
      
type alias Model =
    {image : String
    ,gridx : String
    ,gridy : String
    ,text : String
    ,size : String
    ,position: String
    ,captions: List Caption
    }

init : String -> (Model, Cmd Msg)
init image =
    ( { image = image
      , gridx = "0"
      , gridy = "0"
      , text = ""
      , size = ""
      , position = ""
      , captions = [] }
    , Cmd.none
    )

-- UPDATE

type Msg
    = Grid
    | AddCap
    | GotImg (Result Http.Error String)
    | ChangeX String
    | ChangeY String
    | ChangeText String
    | ChangeSize String
    | ChangePosition String

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        AddCap ->
            if model.text == "" ||
                model.size == "" ||
                    model.position == "" then
                (model, Cmd.none )
                    else
                        ( { model | captions =
                                {text = model.text
                                , size = model.size
                                , position = model.position} :: model.captions },
                              Cmd.none )
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
            ( { model | image = model.image }, addGrid model)
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
                , style "width" (String.fromInt (2*size) ++ "px"), onInput ChangeX ] []
        , input [ placeholder "grid v"
                , style "width" (String.fromInt (2*size) ++ "px"), onInput ChangeY ] []
        , button [ onClick Grid ] [ text "grid" ]
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
                , style "width" (String.fromInt (4*size//3) ++ "px")
                , onInput ChangeText ] []
        , input [ placeholder "Size"
                , style "width" (String.fromInt (4*size//3) ++ "px")
                , onInput ChangeSize ] []
        , input [ placeholder "Position"
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
             ++ List.map toRow model.captions
        )
    ]
      
-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none

-- HTTP
addGrid : Model -> Cmd Msg
addGrid model =
    Http.get
        { url = "http://localhost:8080/grid/?grid=" ++ model.gridx ++ "," ++ model.gridy
        , expect = Http.expectString GotImg
        }

-- addCaption : Model -> 
