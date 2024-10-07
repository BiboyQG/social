import { confirm } from '../utils'
import PropTypes from 'prop-types'
import { useNavigate } from 'react-router-dom'
import { useState } from 'react'

export default function Button({ token }) {
    const navigate = useNavigate()
    const [isLoading, setIsLoading] = useState(false)

    const handleConfirm = async () => {
        setIsLoading(true)
        try {
            const status = await confirm(token)
            if (status === 201) {
                navigate('/success')
            } else {
                navigate('/error')
            }
        } catch (error) {
            console.error('Error confirming account:', error)
            navigate('/error')
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <button
            className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 disabled:opacity-50 disabled:cursor-not-allowed"
            onClick={handleConfirm}
            disabled={isLoading}
        >
            {isLoading ? (
                <span className="flex items-center">
                    <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Confirming...
                </span>
            ) : (
                'Confirm Account'
            )}
        </button>
    );
}

Button.propTypes = {
    token: PropTypes.string.isRequired,
}