import {
  ADD_RECORD,
  SAVE_RECORD
} from '../actions'

export default function records(state = [], action){
  switch (action.type) {
    case ADD_RECORD:
      return [
        ...state,
        {
          id: action.id,
          name: action.name,
          saved: false
        }
      ]
    case SAVE_RECORD:
      return state.map((record, index) => {
        if (index === action.id) {
          return Object.assign({}, record, {
            saved: true
          })
        }
        return record
      })
    default:
      return state
  }
}
