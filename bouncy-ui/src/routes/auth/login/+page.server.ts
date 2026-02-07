import axios from "axios";
import {setAuthCookie} from "$lib/services/cookie_utils";
import {fail, redirect} from "@sveltejs/kit";

type AuthError = {
    success: false;
    error:
        | "invalid-input"
        | "invalid-credentials"
        | "connection-error"
        | "database-error"
        | "login-failed"
    message: string;
}

type AuthSuccess = {
    success: true;
};

type AuthResult = AuthError | AuthSuccess;


export const actions = {
    default: async ({cookies, request}) => {
        const data = await request.formData()
        const email = data?.get('email')?.toString() || ""
        const password = data?.get('password')?.toString() || ""

        const authenticateUser = async (): Promise<AuthResult> => {
            try {
                if (!email || !password) {
                    return {
                        success: false,
                        error: "invalid-input",
                        message: "Email and Password are required"
                    }
                }

                const response = await axios.post("http://localhost:3000/api/v1/auth/login",
                    {email, password})

                if (response.status !== 200) {
                    return {
                        success: false,
                        error: "login-failed",
                        message: "Authenticated successfully",
                    }
                }

                setAuthCookie(cookies, response.data.token)

                return {success: true};
            } catch (error) {
                // Get error message from any type of error
                const errorMessage =
                    error instanceof Error ? error.message : String(error);

                // Simple error classification based on key terms
                let errorType: AuthError["error"] = "login-failed";
                let errorMsg = "An unexpected error occurred";

                // Simple keyword-based error detection
                if (
                    errorMessage.includes("network") ||
                    errorMessage.includes("connect")
                ) {
                    errorType = "connection-error";
                    errorMsg =
                        "Unable to connect to the service. Please try again later.";
                } else if (
                    errorMessage.includes("database") ||
                    errorMessage.includes("query")
                ) {
                    errorType = "database-error";
                    errorMsg = "Database error. Please try again later.";
                }
                console.warn(errorType, errorMessage);
                return {
                    success: false,
                    error: errorType,
                    message: errorMsg,
                };
            }
        }
        const result = await authenticateUser()

        if (!result.success) {
            return handleError(result);
        }

        // Login succeeded, perform redirect
        console.log("Login successful, redirecting to dashboard");
        throw redirect(302, "/");

    }
}

function handleError(result: AuthError): ReturnType<typeof fail> {
    // Simple mapping of error types to status codes
    let statusCode = 500;

    // Define possible response shapes
    type ErrorResponse = { error: boolean; message: string };
    type CredentialsResponse = {
        credentials: boolean;
        message: string;
    };
    type InvalidResponse = { invalid: boolean; message: string };

    // Start with default error response
    let responseData:
        | ErrorResponse
        | CredentialsResponse
        | InvalidResponse = {error: true, message: result.message};

    if (result.error === "invalid-credentials") {
        statusCode = 400;
        responseData = {credentials: true, message: result.message};
    } else if (result.error === "invalid-input") {
        statusCode = 400;
        responseData = {invalid: true, message: result.message};
    } else if (result.error === "connection-error") {
        statusCode = 503;
    }

    return fail(statusCode, responseData);
}