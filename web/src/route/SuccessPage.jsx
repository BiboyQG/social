export default function SuccessPage() {
    return (
        <div className="min-h-screen flex items-center justify-center">
            <div className="bg-white p-8 rounded-lg shadow-md text-center">
                <svg className="w-16 h-16 text-green-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path>
                </svg>
                <h1 className="text-3xl font-bold text-gray-800 mb-4">Account Validated Successfully!</h1>
                <p className="text-gray-600">Thank you for confirming your account. You can now enjoy all the features of our platform.</p>
            </div>
        </div>
    )
}