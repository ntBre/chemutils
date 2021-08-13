module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
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

type alias Model = { image : String }

init : String -> (Model, Cmd Msg)
init image =
    ( { image = image }
    , Cmd.none
    )

-- UPDATE

type Msg
    = Grid
    | GotImg (Result Http.Error String)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        Grid ->
            ( { model | image = model.image }, addGrid )
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
    , button [ onClick Grid ] [ text "grid" ]
    , button [] [ text "add caption" ]
    ]

-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none


-- HTTP
addGrid : Cmd Msg
addGrid =
    Http.get
        { url = "http://localhost:8080/grid/?grid=5,5"
        , expect = Http.expectString GotImg
        }
