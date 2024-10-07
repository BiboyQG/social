import axios from "axios";

export function confirm(token) {
    return axios.put(`/users/activate/${token}`)
        .then(response => {
            console.log(response.status);
            return response.status;
        })
        .catch(error => {
            console.log(error.response);
            return error.response.status;
        });
}