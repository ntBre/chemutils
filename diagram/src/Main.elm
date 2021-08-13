module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)

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

type Msg = NoOp

update : Msg -> Model -> (Model, Cmd Msg)
update _ model =
    ( model, Cmd.none )

-- VIEW

view : Model -> Html Msg
view model =
  div []
    [ img [src model.image] []
    , button [] [ text "grid" ]
    , button [] [ text "add caption" ]
    ]

-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none
