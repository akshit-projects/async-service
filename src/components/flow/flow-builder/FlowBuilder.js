import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Button, Alert, TextField, Grid } from "@mui/material";
import Modal from "react-modal";
import "./FlowBuilder.css";
import FlowStep from "./FlowStep";
import { v4 as uuidv4 } from "uuid";
import axios from "axios";
import StepResponse from "./step-response/StepResponse";
import constants from "../../../constants/constants";

Modal.setAppElement("#root");

const WorkflowBuilder = () => {
  const location = useLocation();
  const [steps, setSteps] = useState([]);
  const [selectedStep, setSelectedStep] = useState(null);
  const [error, setError] = useState(false);
  const [flowName, setFlowName] = useState(constants.DEFAULT_FLOW_NAME);
  const [disableFlowActions, setDisableFlowActions] = useState(false);
  const navigate = useNavigate();
  const [isUpdate, setIsUpdate] = useState(false);

  useEffect(() => {
    const pathArr = location.pathname.split("/");
    if (pathArr[2] !== constants.FLOW_NEW_PATH_SUFFIX) {
      const flowId = pathArr[2];
      setIsUpdate(true);
      const options = {
        url: `${constants.BACKEND_URL}/api/v1/flow/${flowId}`,
      };
      axios(options)
        .then((resp) => {
          const data = resp.data;
          const steps = data.steps;
          setFlowName(data.name);
          const flowSteps = steps.map((step) => {
            return {
              functionType: step.function,
              id: step.id || uuidv4(),
              state: {},
              value: {
                [step.function]: step.meta,
              },
            };
          });
          console.log(flowSteps);
          setSteps(flowSteps);
        })
        .catch((err) => {
          setError("Invaild flow id");
          setDisableFlowActions(true);
        });
    }
  }, [location]);

  const addStep = () => {
    setSteps([...steps, { functionType: "", id: uuidv4(), state: {} }]);
    setSelectedStep(steps.length);
  };

  const deleteStep = (index) => {
    const updatedSteps = steps.filter((_, i) => i !== index);
    setSteps(updatedSteps);
    closeModal();
  };

  const updateStep = (index, field, value) => {
    const updatedSteps = [...steps];
    updatedSteps[index][field] = value;
    setSteps(updatedSteps);
  };

  const openModal = (index) => {
    setSelectedStep(index);
  };

  const closeModal = () => {
    setSelectedStep(null);
  };

  const resetStepsStatus = () => {
    steps.forEach((step, idx) => {
      console.log(step);
      if (step.state.status === constants.FLOW_RESPONSE_STATES.PROGRESS)
        updateStep(idx, "state", {});
    });
  };

  const getFlowBody = () => {
    const stepsPayload = steps.map((step, idx) => {
      updateStep(idx, "state", {
        status: constants.FLOW_RESPONSE_STATES.PROGRESS,
      });
      return {
        name: `Step: ${idx + 1}`,
        function: step.functionType,
        meta: step.value[step.functionType],
        id: step.id || uuidv4(),
      };
    });
    const body = JSON.stringify({
      name: flowName,
      steps: stepsPayload,
    });

    return body;
  };

  const saveFlow = () => {
    resetStepsStatus();
    setDisableFlowActions(true);
    const body = getFlowBody();
    console.log(body);
    const options = {
      url: `${constants.BACKEND_URL}/api/v1/flow`,
      method: "POST",
      data: body,
      headers: {
        "Content-Type": "application/json",
      },
    };
    axios(options)
      .then((res) => {
        console.log(res.data);
        navigate(constants.PATHS.FLOWS);
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => {
        setDisableFlowActions(false);
      });
  };

  const updateFlow = () => {
    resetStepsStatus();
    setDisableFlowActions(true);
    const body = getFlowBody();
    console.log(body);
    const options = {
      url: `${constants.BACKEND_URL}/api/v1/flow`,
      method: "PUT",
      data: body,
      headers: {
        "Content-Type": "application/json",
      },
    };
    axios(options)
      .then((res) => {
        console.log(res.data);
        navigate(constants.PATHS.FLOWS);
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => {
        setDisableFlowActions(false);
      });
  };

  const runTests = (e) => {
    setError();
    setDisableFlowActions(true);
    const body = getFlowBody();
    const socket = new WebSocket(`${constants.WS_BACKEND_URL}/api/v1/flow/run`);
    socket.onopen = () => {
      socket.send(body); // send the body
    };
    socket.onmessage = (msg) => {
      try {
        const data = JSON.parse(msg.data);
        // most probably validation error
        if (data.status === constants.FLOW_RESPONSE_STATES.ERROR && !data.id) {
          setError(data.response.error);
          setDisableFlowActions(false);
          socket.close();
          return;
        } else if (data.status === constants.FLOW_RESPONSE_STATES.ERROR) {
          // id is present, so this is an step
          const stepIdx = steps.findIndex((step) => {
            return step.id === data.id;
          });
          if (stepIdx !== -1) {
            updateStep(stepIdx, "state", {
              status: constants.FLOW_RESPONSE_STATES.ERROR,
              response: data.response,
            });
          } else {
            setError("Invalid state, please try again.");
          }
          socket.close();
        } else {
          // a successful step
          const stepIdx = steps.findIndex((step) => {
            return step.id === data.id;
          });
          if (stepIdx !== -1) {
            updateStep(stepIdx, "state", {
              status: constants.FLOW_RESPONSE_STATES.SUCCESS,
              response: data.response,
            });
          } else {
            setError("Invalid state, please try again.");
          }
        }
      } catch (err) {
        setError("An unknown error occurred. Please contact admin.");
      }
    };

    socket.onclose = () => {
      resetStepsStatus();
      setDisableFlowActions(false);
    };
  };

  return (
    <div className="workflow-builder">
      <h2>Flow Builder</h2>
      <div className="workflow-header">
        <TextField
          label="Flow name"
          value={flowName}
          type="text"
          onChange={(e) => setFlowName(e.target.value)}
        />
        <Button
          className="add-step-button"
          variant="contained"
          color="primary"
          disabled={disableFlowActions}
          onClick={addStep}
        >
          Add Step
        </Button>
      </div>
      {error && (
        <Alert sx={{ margin: "1em 0", width: "100%" }} severity="error">
          {error}
        </Alert>
      )}
      <div className="steps">
        {steps.map((step, index) => (
          <div className="step-row">
            <div onClick={(e) => openModal(index)}>
              <span className="step-index">Step {index + 1}</span>
              <span className="step-function">
                {step.functionType || "Unset"}
              </span>
              <StepResponse state={step.state} />
            </div>
            <Button color="primary" onClick={() => deleteStep(index)}>
              <i className="material-icons">delete</i>
            </Button>
          </div>
        ))}
      </div>
      {steps.length ? (
        <Grid container sx={{ marginLeft: "auto", width: "100%" }}>
          <Grid item xs={8}></Grid>
          <Grid item xs={2}>
            <Button
              className="run-flow-button"
              variant="contained"
              onClick={runTests}
              disabled={disableFlowActions}
            >
              Run Tests
            </Button>
          </Grid>
          <Grid item xs={2}>
            <Button
              variant="contained"
              onClick={isUpdate ? updateFlow : saveFlow}
              disabled={disableFlowActions}
            >
              {isUpdate ? "Upd" : "Save"} Flow
            </Button>
          </Grid>
        </Grid>
      ) : (
        <></>
      )}
      <Modal
        isOpen={selectedStep !== null}
        onRequestClose={closeModal}
        className="modal"
        overlayClassName="overlay"
      >
        {selectedStep !== null && (
          <>
            <div className="modal-header">
              <h3 className="ml-3">Step {selectedStep + 1} Details</h3>
              <Button
                onClick={closeModal}
                sx={{ padding: "0 !important", minWidth: "unset" }}
              >
                <i className="material-icons">close</i>
              </Button>
            </div>
            <FlowStep
              key={selectedStep}
              step={steps[selectedStep]}
              index={selectedStep}
              onDelete={deleteStep}
              onUpdate={updateStep}
            />
          </>
        )}
      </Modal>
    </div>
  );
};

export default WorkflowBuilder;
