import fetch from 'cross-fetch'
export const REQUEST_RECORDS = 'REQUEST_RECORDS'
export const RECEIVE_RECORDS = 'RECEIVE_RECORDS'
export const SELECT_PERIOD = 'SELECT_PERIOD'
export const INVALIDATE_PERIOD = 'INVALIDATE_PERIOD'
export const FETCH_RECORDS_FAIL = 'FETCH_RECORDS_FAIL'

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

export function receiveRecords(period, json){
  return {
    type: RECEIVE_RECORDS,
    period,
    records: json,
    receivedAt: Date.now()
  }
}

function fetchRecords(period){
  return dispatch => {
    dispatch(requestRecords(period))

    const userID = 42
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
    const timestamp = Math.floor(Date.now() / 1000)
    const url = `http://localhost/time-tracker/records?user_id=${userID}&ts=${timestamp}&tz=${timezone}&period=${period}`
    return fetch(url)
      // Try to parse the response
      .then(response =>
        response.json().then(json => ({
          status: response.status,
          json
        })
      ))
      .then(
        // Both fetching and parsing succeeded
        ({ status, json }) => {
          if (status >= 400) {
            dispatch({type: FETCH_RECORDS_FAIL, err: "error: status code "+status, period: period, receivedAt: Date.now()})
          } else {
            dispatch(receiveRecords(period, json))
          }
        },
        // Either fetching or parsing failed!
        err => {
          dispatch({type: FETCH_RECORDS_FAIL, err: err.message, period: period, receivedAt: Date.now(), records: []})
        }
      )
    }
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

