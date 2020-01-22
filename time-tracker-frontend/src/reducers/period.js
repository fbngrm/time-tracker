import {
  SELECT_PERIOD,
  INVALIDATE_PERIOD,
  REQUEST_RECORDS,
  RECEIVE_RECORDS,
  SAVE_RECORD
} from '../actions'

export function selectedPeriod(state = 'day', action) {
  switch (action.type) {
    case SELECT_PERIOD:
      return action.period
    default:
      return state
  }
}

function records(
  state = {
    isFetching: false,
    didInvalidate: false,
    items: []
  },
  action
) {
  switch (action.type) {
    case INVALIDATE_PERIOD:
      return Object.assign({}, state, {
        didInvalidate: true
      })
    case REQUEST_RECORDS:
      return Object.assign({}, state, {
        isFetching: true,
        didInvalidate: false
      })
    case RECEIVE_RECORDS:
      return Object.assign({}, state, {
        isFetching: false,
        didInvalidate: false,
        items: action.records,
        lastUpdated: action.receivedAt
      })
    case SAVE_RECORD:
      return Object.assign({}, state, {
        isFetching: false,
        didInvalidate: false,
        items: [action.record, ...state.items],
        lastUpdated: action.receivedAt
      })
    default:
      return state
  }
}

export function recordsByPeriod(state = {}, action) {
  switch (action.type) {
    case INVALIDATE_PERIOD:
    case RECEIVE_RECORDS:
    case REQUEST_RECORDS:
      return Object.assign({}, state, {
        [action.period]: records(state[action.period], action)
      })
    case SAVE_RECORD:
      for (const period in state) {
         state = Object.assign({}, state, {
             [period]: records(state[period], action)
           })
         }
        return state
    default:
      return state
  }
}
