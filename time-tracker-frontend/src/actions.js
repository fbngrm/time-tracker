import fetch from 'cross-fetch'
export const REQUEST_RECORDS = 'REQUEST_RECORDS'
export const RECEIVE_RECORDS = 'RECEIVE_RECORDS'
export const SELECT_PERIOD = 'SELECT_PERIOD'
export const INVALIDATE_PERIOD = 'INVALIDATE_PERIOD'

export function selectPeriod(period) {
  return {
    type: SELECT_PERIOD,
    period
  }
}

export function invalidatePeriod(period) {
  return {
    type: INVALIDATE_PERIOD,
    period
  }
}

function requestRecords(period) {
  return {
    type: REQUEST_RECORDS,
    period
  }
}

export const receiveRecords = (period, json) => ({
  type: RECEIVE_RECORDS,
  period,
  records: json,
  receivedAt: Date.now()
})

const fetchRecords = period => dispatch => {
  dispatch(requestRecords(period))

  const userID = 42
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
  if (!Date.now) {
    Date.now = function() { return new Date().getTime() }
  }
  const timestamp = Math.floor(Date.now() / 1000)
  const url = `http://localhost:8080/records?user_id=${userID}&ts=${timestamp}&tz=${timezone}&period=${period}`
    console.log(url)
  return fetch(url)
    .then(response => response.json())
    .then(json => dispatch(receiveRecords(period, json)))
}

function shouldFetchRecords(state, period) {
  const records = state.recordsByPeriod[period]
  if (!records) {
    return true
  } else if (records.isFetching) {
    return false
  } else {
    return records.didInvalidate
  }
}

export function fetchRecordsIfNeeded(period) {
  return (dispatch, getState) => {
    if (shouldFetchRecords(getState(), period)) {
      return dispatch(fetchRecords(period))
    }
  }
}
