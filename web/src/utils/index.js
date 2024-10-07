import axios from "axios";

export function confirm(token) {
    return axios.get(`/authentication/user/${token}`)
        .then(response => {
            console.log(response.status);
            return response.status;
        })
        .catch(error => {
            console.log(error.response);
            return error.response.status;
        });
}