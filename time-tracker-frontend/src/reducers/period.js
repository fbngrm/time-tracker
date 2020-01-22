import {
  SELECT_PERIOD,
  INVALIDATE_PERIOD,
  REQUEST_RECORDS,
  RECEIVE_RECORDS
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
    default:
      return state
  }
}
