module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Debug

-- MAIN

main =
  Browser.element
      { init = init
      , view = view
      , update = update
      , subscriptions = subscriptions
      }

-- MODEL

type alias Model =
    {image : String
    ,gridx : String
    ,gridy : String
    }

init : String -> (Model, Cmd Msg)
init image =
    ( { image = image, gridx = "0", gridy = "0" }
    , Cmd.none
    )

-- UPDATE

type Msg
    = Grid
    | GotImg (Result Http.Error String)
    | ChangeX String
    | ChangeY String

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        ChangeX newX ->
            ( { model | gridx = newX }, Cmd.none )
        ChangeY newY ->
            ( { model | gridy = newY }, Cmd.none )
        Grid ->
            ( { model | image = model.image }, addGrid model)
        GotImg result ->
            case result of
                Ok img ->
                    let dum = Debug.log "received img" in
                    ( {model | image = img}, Cmd.none)
                Err _ ->
                    let dum = Debug.log "did not receive img" in
                    (model, Cmd.none)

-- VIEW

view : Model -> Html Msg
view model =
  div []
    [ img [src model.image] []
    , input [ placeholder "grid x", onInput ChangeX ] []
    , input [ placeholder "grid y", onInput ChangeY ] []
    , button [ onClick Grid ] [ text "grid" ]
    , button [] [ text "add caption" ]
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
