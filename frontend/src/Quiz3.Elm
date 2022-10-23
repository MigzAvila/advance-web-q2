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
