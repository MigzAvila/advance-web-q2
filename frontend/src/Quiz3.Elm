module Quiz3 exposing (main)

import Browser
-- Bring everything from Html Like import namespace
import Html exposing (..)
import Html.Attributes exposing (..)
-- Import Events that will be used
import Html.Events exposing (onClick, onInput, onSubmit)

-- Create a Type Todo with id, name(title), and complete_status Single Todo
type alias Todo =
    { id : Int
    , name : String
    , isComplete: Bool
    }

-- Create a Type Model to wrap all the todos (Holds multiple todo)
type alias Model =
    { field : String
    , uid : Int
    , todos : List Todo
    }

-- Create a Type Msg for the displayed todo - Completed, Name, Delete Todo 
type Msg
    = AddTodo
    | SetField String
    | DeleteTodo Int
    | CompleteTodo Int Bool


-- create and Initialize Model to be empty {"", 0, []}
initialModel : Model
initialModel =
    { field = ""
    , uid = 0
    , todos = []
    }


-- Create model
view : Model -> Html Msg
view model =
    div [ class "text-center" ] [ div [ class "todo-container text-left p-24 bg-white shadow-sm rounded flex flex-col mx-auto my-48" ]
        [ header [ ]
            [ h1 [ class "title" ] [ text "Todo List" ] ]
        , Html.form [ class "form" ,onSubmit AddTodo ] [
            input
                [ class "todo-input"
                , placeholder "Title here.."
                , onInput
                    (\string -> SetField string)
                , value model.field
                ]
                []
            , button [ class "btn", type_ "submit", disabled (model.field == "") ] [ text "Create" ]
        ]
        , ul [ class "text-left mt-24" ] (List.map viewSearchResult model.todos)
    ]
    ]
