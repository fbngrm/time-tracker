import fetch from 'cross-fetch'

export const REQUEST_RECORDS = 'REQUEST_RECORDS'
export const RECEIVE_RECORDS = 'RECEIVE_RECORDS'
export const ADD_RECORD = 'ADD_RECORD'
export const SAVE_RECORD = 'SAVE_RECORD'
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
    const url = `http://localhost:8080/records?user_id=${userID}&ts=${timestamp}&tz=${timezone}&period=${period}`
    return fetch(url)
      .then(response => response.json())
      .then(json => dispatch(receiveRecords(period, json)))
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

let nextRecordId = 0
export function addRecord(name) {
  return {
    type: ADD_RECORD,
    id: nextRecordId++,
    name
  }
}

function savedRecord(id, json){
  return {
    type: SAVE_RECORD,
    record: json,
    receivedAt: Date.now(),
    id
  }
}

function saveRecord(record){
  return dispatch => {
    // dispatch(requestRecords(period))
    const { id, name, start, startLoc, stop, stopLoc, duration } = record
    const url = `http://localhost:8080/record`
    return fetch(url,
       {
        method: 'POST',
        headers:{'Content-Type': 'application/json'},
        body: JSON.stringify({
          user_id: 42,
          name: name,
          start_time: start,
          start_loc: startLoc,
          stop_time: stop,
          stop_loc: stopLoc,
          duration: duration
        })
      })
      .then(response => response.json())
      .then(json => dispatch(savedRecord(id, json)))
  }
}

function shouldSaveRecord(record) {
  if (record.isFetching) {
    return false
  }
  return !record.saved
}

export function saveRecordIfNeeded(record) {
  return (dispatch, getState) => {
    if (shouldSaveRecord(getState(), record)) {
      return dispatch(saveRecord(record))
    }
  }
}
