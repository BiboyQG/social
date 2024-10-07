import Button from '../component/Button'
import { useParams } from 'react-router-dom'

export default function ConfirmPage() {
    const { token } = useParams()

	return (
		<div>
			<h1>Confirm Page</h1>
			<Button token={token}/>
		</div>
	);
}
