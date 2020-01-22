import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Record from './Record'

export default class RecordList extends Component {
  render() {
    const { records, saveRecord } = this.props
    return (
      <div>
        {records.map(record =>
          <Record
            {...record}
            key={record.id}
            onClick={saveRecord}
          />
        )}
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
