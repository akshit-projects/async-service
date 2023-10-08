import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { login } from "./auth-utils";
import constants from "../../constants/constants";
import axios from "axios";
import './Login.css';
const Login = () => {
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const token = localStorage.getItem('token');

  useEffect(() => {
    const sessionId = localStorage.getItem("token");
    if (sessionId) {
      navigate("/");
    }
  }, [navigate]);

  window.handleCredentialResponse = async (response) => {
    const idToken = response.credential;

    try {
        const options = {
            url: `${constants.BACKEND_URL}/api/v1/login`,
            method: 'POST',
            data: { idToken },
            credentials: 'include',
        };
        const response = await axios(options);
        const data = response.data;
        if (!data.token) {
            throw new Error('Unable to authenticate. Please try again later.');
        }
        login(data);
        navigate("/flow");
    } catch (error) {
        if (error instanceof TypeError) {
            setError(
                'An error occurred while logging in. Please try again later or contact AkshitBansal.'
            );
        } else {
            setError(
                error.message,
            );
        }
    }
  };

  return (
    <div className="center-page">
        {!token &&
            <div>
                <h2 style={{textAlign: 'center'}}>Login With Google</h2>
                <div
                    id="g_id_onload"
                    data-client_id="657098499628-hlccv6fl77gbj4srer1asb5v66l1chi6.apps.googleusercontent.com"
                    data-callback="handleCredentialResponse"
                ></div>
                <div className="g_id_signin w-full" style={{ textAlign: 'center' }} data-logo_alignment="center" data-type="standard"></div>
                {error && <p>{error}</p>}
            </div>
        }
    </div>
  );
};

export default Login;
