import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import {
  selectPeriod,
  fetchRecordsIfNeeded,
  invalidatePeriod
} from '../actions'
import Picker from '../components/Picker'
import Records from '../components/Records'

class AsyncApp extends Component {
  constructor(props) {
    super(props)
    this.handleChange = this.handleChange.bind(this)
    this.handleRefreshClick = this.handleRefreshClick.bind(this)
  }

  componentDidMount() {
    const { dispatch, selectedPeriod } = this.props
    dispatch(fetchRecordsIfNeeded(selectedPeriod))
  }

  componentDidUpdate(prevProps) {
    if (this.props.selectedPeriod !== prevProps.selectedPeriod) {
      const { dispatch, selectedPeriod } = this.props
      dispatch(fetchRecordsIfNeeded(selectedPeriod))
    }
  }

  handleChange(nextPeriod) {
    this.props.dispatch(selectPeriod(nextPeriod))
    this.props.dispatch(fetchRecordsIfNeeded(nextPeriod))
  }

  handleRefreshClick(e) {
    e.preventDefault()
    const { dispatch, selectedPeriod } = this.props
    dispatch(invalidatePeriod(selectedPeriod))
    dispatch(fetchRecordsIfNeeded(selectedPeriod))
  }

  render() {
    const { selectedPeriod, records, isFetching, lastUpdated } = this.props
    return (
      <div>
        <Picker
          value={selectedPeriod}
          onChange={this.handleChange}
          options={['day', 'week', 'month']}
        />
        <p>
          {lastUpdated && (
            <span>
              Last updated at {new Date(lastUpdated).toLocaleTimeString()}.{' '}
            </span>
          )}
          {!isFetching && (
            <button onClick={this.handleRefreshClick}>Refresh</button>
          )}
        </p>
        {isFetching && records.length === 0 && <h2>Loading...</h2>}
        {!isFetching && records.length === 0 && <h2>Empty.</h2>}
        {records.length > 0 && (
          <div style={{ opacity: isFetching ? 0.5 : 1 }}>
            <Records records={records} />
          </div>
        )}
      </div>
    )
  }
}

AsyncApp.propTypes = {
  selectedPeriod: PropTypes.string.isRequired,
  records: PropTypes.array.isRequired,
  isFetching: PropTypes.bool.isRequired,
  lastUpdated: PropTypes.number,
  dispatch: PropTypes.func.isRequired
}

function mapStateToProps(state) {
  const { selectedPeriod, recordsByPeriod } = state
  const { isFetching, lastUpdated, items: records } = recordsByPeriod[
    selectedPeriod
  ] || {
    isFetching: true,
    items: []
  }
  return {
    selectedPeriod,
    records,
    isFetching,
    lastUpdated
  }
}
export default connect(mapStateToProps)(AsyncApp)
