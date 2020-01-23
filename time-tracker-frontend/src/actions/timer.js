import fetch from 'cross-fetch'
export const ADD_TIMER = 'ADD_TIMER'
export const START_TIMER = 'START_TIMER'
export const STOP_TIMER = 'STOP_TIMER'
export const SAVE_TIMER = 'SAVE_TIMER'
export const RESET_TIMER = 'RESET_TIMER'
export const SAVE_RECORD = 'SAVE_RECORD'

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
    start: Math.floor(Date.now() / 1000),
    startedAt: (state.startedAt === -1) ? Math.floor(Date.now() / 1000) : state.startedAt, // set initial start time once
    startLoc: (state.startLoc === "") ? Intl.DateTimeFormat().resolvedOptions().timeZone : state.startLoc, // set initial start time once
    isRunning: true,
    isStopped: false
  }
}

export function stopTimer(state){
  const now = Math.floor(Date.now() / 1000)
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
    const url = `http://localhost:8080/record`
    return fetch(url,
       {
        method: 'POST',
        headers:{'Content-Type': 'application/json'},
        body: JSON.stringify({
          user_id: 42,
          name: name,
          start_time: startedAt, // seconds since unix epoch,
          start_loc: startLoc,
          stop_time :stoppedAt, // seconds since unix epoch,
          stop_loc: stopLoc,
          duration:time // seconds
        })
      })
      .then(response => response.json())
      .then(json => dispatch(saveRecord(json)))
      .then(dispatch(resetTimer()))
  }
}
