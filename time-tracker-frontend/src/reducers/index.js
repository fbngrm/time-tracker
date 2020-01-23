import { combineReducers } from 'redux'
import { recordsByPeriod, selectedPeriod } from './period'
import timer from './timer'

const rootReducer = combineReducers({
  recordsByPeriod,
  selectedPeriod,
  timer
})

export default rootReducer
