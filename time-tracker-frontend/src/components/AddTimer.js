import React, { Component } from 'react'
import PropTypes from 'prop-types'

export default class AddTimer extends Component {
  constructor(props){
    super(props)
    this.state = {
      start: 0,
      time: 0
    }
    this.startTimer = this.startTimer.bind(this)
    this.stopTimer = this.stopTimer.bind(this)
    this.saveTimer = this.saveTimer.bind(this)
  }

  startTimer() {
    this.setState({
      time: this.state.time, // timer
      start: Math.floor(Date.now() / 1000) - this.state.time // set every time the timer is started/continued
    })
    this.timer = setInterval(() => this.setState({
      time: Math.floor(Date.now() / 1000) - this.state.start
    }), 1000)
    this.props.startTimer(this.props.timer)
  }

  stopTimer() {
    this.props.stopTimer(this.props.timer)
    clearInterval(this.timer)
  }

  saveTimer() {
    this.setState({
      time: 0,
      start: 0
    })
    this.props.saveTimer(this.props.timer)
  }

  render() {
    const { addTimer, saveTimer, timer } = this.props

    let input

    let startButton = (!timer.isRunning && timer.startedAt !== undefined) ?
      <button onClick={this.startTimer}>start</button> :
      <button disabled="disabled">start</button>

    let stopButton = (timer.isRunning) ?
      <button onClick={this.stopTimer}>stop</button> :
      <button disabled="disabled">stop</button> 

    let saveButton = (timer.isStopped) ?
      <button onClick={this.saveTimer}>save</button> :
      <button disabled="disabled">save</button>

    let visibleTimer = (timer.startedAt === undefined) ? 'hide' : ''

    return (
      <div>
        <div>
          <form onSubmit={e => {
            e.preventDefault()
            if (!input.value.trim()) {
              return
            }
            addTimer(input.value)
            input.value = ''
          }}>
            <input ref={node => input = node} />
            <button type="submit">
              add
            </button>
            {startButton}
            {stopButton}
            {saveButton}
          </form>
        </div>
        <p className={visibleTimer}>
          <span>{timer.name}</span>
          <span>{formatTime(timer.startedAt)}</span>
          <span>{formatTime(timer.stoppedAt)}</span>
          <span>{formatTimer(this.state.time)}</span>
        </p>
      </div>
    )
  }
}

function formatTimer(t) {
  if (t === 0) return ""
  var sec_num = parseInt(t, 10)
  var hours   = Math.floor(sec_num / 3600)
  var minutes = Math.floor((sec_num - (hours * 3600)) / 60)
  var seconds = sec_num - (hours * 3600) - (minutes * 60)
  if (hours   < 10) {hours   = "0"+hours}
  if (minutes < 10) {minutes = "0"+minutes}
  if (seconds < 10) {seconds = "0"+seconds}
  return hours+':'+minutes+':'+seconds
}

function formatTime(t) {
  if (t === -1) return ""
  var a = new Date(t*1000);
  var months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
  var year = a.getFullYear()
  var month = months[a.getMonth()]
  var date = a.getDate()
  var hour = a.getHours()
  var min = a.getMinutes() < 10 ? '0' + a.getMinutes() : a.getMinutes();
  var sec = a.getSeconds() < 10 ? '0' + a.getSeconds() : a.getSeconds()
  var time = date + ' ' + month + ' ' + year + ' ' + hour + ':' + min + ':' + sec
  return time
}

AddTimer.propTypes = {
  timer: PropTypes.shape({
  }).isRequired
}
