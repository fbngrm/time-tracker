import { connect } from 'react-redux'
import { saveRecordIfNeeded } from '../actions'
import Timer from '../components/Timer'

function mapStateToProps(state){
  return {
    record: state.record
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
)(Timer)
