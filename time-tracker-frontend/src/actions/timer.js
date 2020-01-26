import fetch from 'cross-fetch'
export const ADD_TIMER = 'ADD_TIMER'
export const START_TIMER = 'START_TIMER'
export const STOP_TIMER = 'STOP_TIMER'
export const SAVE_TIMER = 'SAVE_TIMER'
export const RESET_TIMER = 'RESET_TIMER'
export const SAVE_RECORD = 'SAVE_RECORD'
export const SAVE_TIMER_FAIL = 'SAVE_TIMER_FAIL'

export function addTimer(name) {
  return {
    type: ADD_TIMER,
    name: name,
    time: 0,
    start: 0,
    startedAt: -1,
    startLoc: ""
  }
}

export function startTimer(state){
  return {
    type: START_TIMER,
    time: state.time,
    start: Date.now(),
    startedAt: (state.startedAt === -1) ? Date.now() : state.startedAt, // set initial start time once
    startLoc: (state.startLoc === "") ? Intl.DateTimeFormat().resolvedOptions().timeZone : state.startLoc, // set initial start time once
    isRunning: true,
    isStopped: false
  }
}

export function stopTimer(state){
  const now = Date.now()
  return {
    type: STOP_TIMER,
    isRunning: false,
    isStopped: true,
    stoppedAt: now,
    stopLoc: Intl.DateTimeFormat().resolvedOptions().timeZone,
    time: state.time + (now - state.start)
  }
}

export function saveTimer(state){
  return {
    type: SAVE_TIMER,
    isSaving: true
  }
}

export function resetTimer(){
  return {
    type: RESET_TIMER
  }
}

export function saveTimerIfNeeded(state) {
  return (dispatch) => {
    if (!state.isSaving) {
      return dispatch(save(state))
    }
  }
}

function saveRecord(json){
  return {
    type: SAVE_RECORD,
    record: json,
    receivedAt: Date.now()
  }
}

function save(state){
  return dispatch => {
    dispatch(saveTimer)
    const { name, startedAt, startLoc, stoppedAt, stopLoc, time } = state
    const url = `http://localhost/time-tracker/record`
    return fetch(url,
       {
        method: 'POST',
        headers:{'Content-Type': 'application/json'},
        body: JSON.stringify({
          user_id: 42,
          name: name,
          start_time: Math.floor(startedAt / 1000), // seconds since unix epoch,
          start_loc: startLoc,
          stop_time: Math.floor(stoppedAt / 1000), // seconds since unix epoch,
          stop_loc: stopLoc,
          duration: Math.floor(time / 1000) // seconds
        })
      })
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
            dispatch({type: SAVE_TIMER_FAIL, err: "error: status code "+status})
          } else {
            dispatch(saveRecord(json))
            dispatch(resetTimer())
          }
        },
        // Either fetching or parsing failed!
        err => {
          dispatch({type: SAVE_TIMER_FAIL, err: err.message})
        }
      )
    }
}
