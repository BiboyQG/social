import Button from '../component/Button'
import { useParams } from 'react-router-dom'

export default function ConfirmPage() {
    const { token } = useParams()

	return (
		<div className="min-h-screen flex items-center justify-center">
			<div className="bg-white p-8 rounded-lg shadow-md text-center">
				<h1 className="text-3xl font-bold mb-6 text-gray-800">Confirm Your Account</h1>
				<p className="mb-6 text-gray-600">Please click the button below to confirm your account.</p>
				<Button token={token} />
			</div>
		</div>
	);
}
