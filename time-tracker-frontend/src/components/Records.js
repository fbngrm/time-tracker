import React, { Component } from 'react'
import PropTypes from 'prop-types'

export default class Records extends Component {
  render() {
    return (
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Start</th>
            <th>Start Location</th>
            <th>Stop</th>
            <th>Stop Location</th>
            <th>Duration</th>
          </tr>
        </thead>
        <tbody>
        {this.props.records.map((record, i) => (
          <tr key={record.record_id}>
            <React.Fragment>
              <td >{record.name}</td>
              <td >{record.start_time}</td>
              <td >{record.start_loc}</td>
              <td >{record.stop_time}</td>
              <td >{record.stop_loc}</td>
              <td >{record.duration}</td>
            </React.Fragment>
          </tr>
        ))}
        </tbody>
      </table>
    )
  }
}
Records.propTypes = {
  records: PropTypes.array.isRequired
}
