import { confirm } from '../utils'
import PropTypes from 'prop-types'
import { useNavigate } from 'react-router-dom'

export default function Button({ token }) {
    const navigate = useNavigate()

	return (
        <button className="btn btn-primary" onClick={() => confirm(token).then(status => {
            if (status === 201) {
                navigate('/success')
            } else {
                navigate('/error')
            }
        })}>
            Confirm
		</button>
	);
}

Button.propTypes = {
    token: PropTypes.string.isRequired,
}