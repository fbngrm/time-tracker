import { combineReducers } from 'redux'
import { recordsByPeriod, selectedPeriod } from './period'
import records from './records'

const rootReducer = combineReducers({
  recordsByPeriod,
  selectedPeriod,
  records
})

export default rootReducer
