import React, { Component } from 'react'
import PropTypes from 'prop-types'

export default class Timer extends Component {
  constructor(props){
    super(props)
    // this.id = props.id
    // this.name = props.name
    // this.state = {
    //   time: 0,
    //   start: 0,
    //   startTime: -1,
    //   startLoc: ""
    // }
    this.startTimer = this.startTimer.bind(this)
    this.stopTimer = this.stopTimer.bind(this)
  }

  startTimer() {
    // var startTime = (this.state.startTime === -1) ? Date.now() : this.state.startTime
    // var startLoc = (this.state.startLoc === "") ? Intl.DateTimeFormat().resolvedOptions().timeZone : this.state.startLoc
    // this.setState({
    //   time: this.state.time, // timer
    //   start: Date.now() - this.state.time, // set every time the timer is started/continued
    //   startTime: startTime, // set once on initial start
    //   startLoc: startLoc,
    //   running: true,
    //   stopped: false,
    // })
    this.timer = setInterval(() => this.setState({
      time: Date.now() - this.state.start
    }), 1)
  }

  stopTimer() {
    // this.setState({
    //   running: false,
    //   stopped: true,
    //   stopTime: Date.now(),
    //   stopLoc: Intl.DateTimeFormat().resolvedOptions().timeZone,
    // })
    clearInterval(this.timer)
  }

  saveTimer() {
    // this.setState({
    //     running: false,
    //     stopped: true,
    //     saved: true,
    // })
    // this.props.saveRecord({
    //     id: this.id,
    //     name: this.name,
    //     start: Math.round(this.state.startTime / 1000), // seconds since unix epoch
    //     startLoc: this.state.startLoc,
    //     stop: Math.round(this.state.stopTime / 1000), // seconds since unix epoch
    //     stopLoc: this.state.stopLoc,
    //     duration: Math.round(this.state.time / 1000),
    // })
  }

  render() {
    // let startButton = (!this.state.running && !this.state.saved)?
    //   <button onClick={this.startTimer}>start</button> :
    //   <button onClick={this.startTimer} disabled="disabled">start</button>
    // let stopButton = (this.state.running && !this.state.saved) ?
    //   <button onClick={this.stopTimer}>stop</button> :
    //   <button onClick={this.stopTimer} disabled="disabled">stop</button> 
    // let saveButton = (this.state.stopped && !this.state.saved )?
    //   <button onClick={this.saveTimer}>save</button> :
    //   <button onClick={this.saveTimer} disabled="disabled">save</button>

    // const { saved } = this.props
    // let savedState = {saved} === true ? 'success' : 'failed'
    const { record } = this.props
    if (record === undefined) {
        return null
    }
    return (
      <div>
      <div className="record">
      <div className="record">
        <span>{record.name}</span>
        <span>{formatTime(this.state.startedAt)}</span>
        <span>{formatTimer(Math.round(this.state.time/1000))}</span>
      </div>
      // <div className="record">
      //   {startButton}
      //   {stopButton}
      //   {saveButton}
      // </div>
      </div>
      </div>
    )
  }
}

function formatTimer(t) {
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
  if (t === -1) {
    return ""
  }
  var a = new Date(t);
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

Timer.propTypes = {
}

