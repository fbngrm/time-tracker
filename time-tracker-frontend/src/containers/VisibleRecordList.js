import { connect } from 'react-redux'
import { saveRecordIfNeeded } from '../actions'
import RecordList from '../components/RecordList'

function mapStateToProps(state){
  return {
    records: state.records
  }
}

function mapDispatchToProps(dispatch){
  return {
    saveRecord: function(state){
      return dispatch(saveRecordIfNeeded(state))
    }
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(RecordList)
