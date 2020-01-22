import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Record from './Record'

export default class RecordList extends Component {
  render() {
    const { records, saveRecord } = this.props
    return (
      <ul>
        {records.map(todo =>
          <Record
            {...todo}
            key={todo.id}
            done={todo.done}
            onClick={saveRecord}
          />
        )}
      </ul>
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
