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
              <td key={record.record_id}>{record.name}</td>
              <td key={record.record_id}>{record.start_time}</td>
              <td key={record.record_id}>{record.start_loc}</td>
              <td key={record.record_id}>{record.stop_time}</td>
              <td key={record.record_id}>{record.stop_loc}</td>
              <td key={record.record_id}>{record.duration}</td>
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
