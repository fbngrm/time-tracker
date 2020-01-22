import React, { Component } from 'react'
import PropTypes from 'prop-types'

export default class Record extends Component {
  constructor(props){
    super(props)
    this.id = props.id
    this.name = props.name
    this.state = {
      time: 0,
      start: 0
    }
    this.startTimer = this.startTimer.bind(this)
    this.stopTimer = this.stopTimer.bind(this)
    this.saveTimer = this.saveTimer.bind(this)
  }

  startTimer() {
    this.setState({
      time: this.state.time,
      start: Date.now() - this.state.time,
      startLoc: Intl.DateTimeFormat().resolvedOptions().timeZone,
      running: true,
      stopped: false,
    })
    this.timer = setInterval(() => this.setState({
      time: Date.now() - this.state.start
    }), 1)
  }

  stopTimer() {
    this.setState({
      running: false,
      stopped: true,
      stopLoc: Intl.DateTimeFormat().resolvedOptions().timeZone,
    })
    clearInterval(this.timer)
  }

  saveTimer() {
    this.setState({
        running: false,
        stopped: true,
        saved: true,
    })
    this.props.onClick({
        id: this.id,
        name: this.name,
        start: Math.round(this.state.start / 1000), // seconds since unix epoch
        startLoc: this.state.startLoc,
        stop: Math.round((this.state.start + this.state.time) / 1000), // seconds since unix epoch
        stopLoc: this.state.stopLoc,
    })
  }

  render() {
    let startButton = (!this.state.running && !this.state.saved)?
      <button onClick={this.startTimer}>start</button> :
      null
    let stopButton = (this.state.running && !this.state.saved) ?
      <button onClick={this.stopTimer}>stop</button> :
      null
    let saveButton = (this.state.stopped && !this.state.saved )?
      <button onClick={this.saveTimer}>save</button> :
      null

    const { saved } = this.props
    let savedState = {saved} ? 'success' : 'failed'

    return (
      <div>
        {this.name}: {formatTime(Math.round(this.state.time/1000))}
        {startButton}
        {stopButton}
        {saveButton}
        {savedState}
      </div>
    )
  }
}

function formatTime(t) {
    var sec_num = parseInt(t, 10)
    var hours   = Math.floor(sec_num / 3600)
    var minutes = Math.floor((sec_num - (hours * 3600)) / 60)
    var seconds = sec_num - (hours * 3600) - (minutes * 60)
    if (hours   < 10) {hours   = "0"+hours}
    if (minutes < 10) {minutes = "0"+minutes}
    if (seconds < 10) {seconds = "0"+seconds}
    return hours+':'+minutes+':'+seconds
}

Record.propTypes = {
  onClick: PropTypes.func.isRequired,
  name: PropTypes.string.isRequired,
  saved: PropTypes.bool.isRequired
}

