import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Button, Alert, TextField, Grid } from "@mui/material";
import Modal from "react-modal";
import "./FlowBuilder.css";
import FlowStep from "./FlowStep";
import { v4 as uuidv4 } from "uuid";
import axios from "axios";
import StepResponse from "./step-response/StepResponse";
import constants, { FLOW_FUNCTIONS } from "../../../constants/constants";
import { validateSteps } from "./util";

Modal.setAppElement("#root");

const WorkflowBuilder = () => {
  const location = useLocation();
  const [steps, setSteps] = useState([]);
  const [selectedStep, setSelectedStep] = useState(null);
  const [error, setError] = useState(false);
  const [flowName, setFlowName] = useState(constants.DEFAULT_FLOW_NAME);
  const [disableFlowActions, setDisableFlowActions] = useState(false);
  const navigate = useNavigate();
  const [flowId, setFlowId] = useState("");

  // if it's an existing flow, get it's data
  useEffect(() => {
    const pathArr = location.pathname.split("/");
    if (pathArr[2] !== constants.FLOW_NEW_PATH_SUFFIX) {
      const flowId = pathArr[2];
      setFlowId(flowId);
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
    } else {
        addStep();
        setSelectedStep(null);
    }
  }, [location]);

  const addStep = () => {
    setSteps([...steps, { functionType: "", id: uuidv4(), state: {}, value: {} }]);
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
      if (step.state.status === constants.FLOW_RESPONSE_STATES.PROGRESS)
        updateStep(idx, "state", {});
    });
  };

  const validateFlow = (flow) => {
    const steps = flow.steps || [];
    const error = validateSteps(steps);
    if (error) {
        setError(error.message);
        return false;
    }

    return true;
  }

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
        timeout: 10000
      };
    });
    const payload = {
      name: flowName,
      steps: stepsPayload,
    };
    if (!validateFlow(payload)) {
        setDisableFlowActions(false);
        resetStepsStatus();
        return null;
    }
    if (flowId) {
      payload.id = flowId;
    }
    const body = JSON.stringify(payload);

    return body;
  };

  const saveFlow = () => {
    resetStepsStatus();
    setDisableFlowActions(true);
    const body = getFlowBody();
    if (!body) {
        return;
    }
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
    if (!body) {
        return;
    }
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
    if (!body) {
        return;
    }
    const socket = new WebSocket(`${constants.WS_BACKEND_URL}/api/v1/flow/run`);
    socket.onopen = () => {
      socket.send(body); // send the body
    };
    socket.onmessage = (msg) => {
      try {
        const data = JSON.parse(msg.data);
        // most probably validation error or a generic error
        if (data.type === constants.FLOW_RESPONSE_STATES.ERROR && !data.stepResponse) {
          setError(data.message);
          setDisableFlowActions(false);
          socket.close();
          return;
        } else if (data.status === constants.FLOW_RESPONSE_STATES.ERROR) {
          // this is a step level error
          const stepResponse = data.stepResponse;
          // id is present, so this is an step
          const stepIdx = steps.findIndex((step) => {
            return step.id === stepResponse.id;
          });
          if (stepIdx !== -1) {
            updateStep(stepIdx, "state", {
              status: constants.FLOW_RESPONSE_STATES.ERROR,
              response: stepResponse,
            });
          } else {
            setError("Invalid state, please try again.");
          }
          socket.close();
        } else {
          const stepResponse = data.stepResponse;  
          // a successful step
          const stepIdx = steps.findIndex((step) => {
            return step.id === stepResponse.id;
          });
          if (stepIdx !== -1) {
            if (stepResponse.status === "ERROR") {
                updateStep(stepIdx, "state", {
                    status: constants.FLOW_RESPONSE_STATES.ERROR,
                    response: stepResponse,
                });
            } else {
                updateStep(stepIdx, "state", {
                    status: constants.FLOW_RESPONSE_STATES.SUCCESS,
                    response: stepResponse,
                });
            }
            
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
      <h1 style={{ alignSelf: 'start', margin: '1.4em !important'}}>Flow Builder</h1>
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
                {step.functionType || "Not Set"}
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
              onClick={flowId ? updateFlow : saveFlow}
              disabled={disableFlowActions}
            >
              {flowId ? "Upd" : "Save"} Flow
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
