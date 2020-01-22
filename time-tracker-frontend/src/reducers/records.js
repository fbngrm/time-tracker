import {
  ADD_RECORD,
  SAVE_RECORD
} from '../actions'

export default function records(state = [], action){
  switch (action.type) {
    case ADD_RECORD:
      return [
        {
          id: action.id,
          name: action.name,
          saved: false
        }
      ]
    default:
      return state
  }
}
