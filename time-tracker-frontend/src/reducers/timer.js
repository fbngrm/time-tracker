import {
  ADD_TIMER,
  START_TIMER,
  STOP_TIMER,
  SAVE_TIMER
} from '../actions'

export default function timer(state = {}, action){
  switch (action.type) {
    case ADD_TIMER:
      return {
          name: action.name,
          time: action.time,
          start: action.start,
          startedAt: action.startedAt,
          startLoc: action.startLoc
        }
    case START_TIMER:
      return Object.assign({}, state, {
        time: action.time,
        start: action.start,
        startedAt: action.startedAt,
        startLoc: action.startLoc,
        isRunning: action.isRunning,
        isStopped: action.isStopped
      })
    case STOP_TIMER:
      return Object.assign({}, state, {
        time: action.time,
        isRunning: action.isRunning,
        isStopped: action.isStopped,
        stoppedAt: action.stoppedAt,
        stopLoc: action.stopLoc
      })
    case SAVE_TIMER:
      return Object.assign({}, state, {
        isSaving: action.isSaving
      })
    default:
      return state
  }
}
