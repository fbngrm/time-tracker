import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Record from './Record'

export default class RecordList extends Component {
  render() {
    const { records, saveRecord } = this.props
    return (
      <div>
      <table>
        <tbody>
        {records.map(record =>
          <tr key={record.record_id}>
            <React.Fragment>
            <Record
              {...record}
              key={record.id}
              onClick={saveRecord}
            />
            </React.Fragment>
          </tr>
        )}
        </tbody>
      </table>
      </div>
    )
  }
}

RecordList.propTypes = {
  records: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    saved: PropTypes.bool.isRequired
  }).isRequired).isRequired,
  saveRecord: PropTypes.func.isRequired
}
