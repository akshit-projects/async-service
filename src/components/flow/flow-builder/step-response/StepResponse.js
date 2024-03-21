import { CircularProgress } from "@material-ui/core";
import React, { useState } from "react";
import SuccessModal from "./SuccessModal";
import ErrorResponse from "./ErrorResponse";
import constants from "../../../../constants/constants";

export default function StepResponse(props) {
  const state = props.state;
  const status = state.status;
  const response = state.response;
  const [isModalOpen, setIsModalOpen] = useState(false);

  const click = (event) => {
    if (event.stopPropagation) {
      event.stopPropagation(); // W3C model
    } else {
      event.cancelBubble = true; // IE model
    }
    setIsModalOpen(true);
  };

  const closeRequest = (event) => {
    if (event.stopPropagation) {
        event.stopPropagation(); // W3C model
      } else {
        event.cancelBubble = true; // IE model
      }
    setIsModalOpen(false);
  }

  let render;
  if (status === constants.FLOW_RESPONSE_STATES.PROGRESS) {
    render = <CircularProgress className="step-status" size={18} />;
  } else if (status === constants.FLOW_RESPONSE_STATES.SUCCESS) {
    render = (
        <>
        <i className="material-icons step-status" onClick={click}>
            check
        </i>
        <SuccessModal isModalOpen={isModalOpen} closeRequest={closeRequest} value={response} />
      </>
    );
  } else if (status === constants.FLOW_RESPONSE_STATES.ERROR) {
    render = <>
        <i className="material-icons step-status close" onClick={click}>close</i>
        <ErrorResponse isModalOpen={isModalOpen} closeRequest={closeRequest} response={response?.response} />
    </>
  } else {
    render = <></>;
  }

  return render;
}
