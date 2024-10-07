export default function ErrorPage() {
    return (
        <div className="min-h-screen flex items-center justify-center">
            <div className="bg-white p-8 rounded-lg shadow-md text-center">
                <svg className="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <h1 className="text-3xl font-bold text-gray-800 mb-4">Account Validation Error</h1>
                <p className="text-gray-600">We encountered an error while validating your account. Please try again or contact support if the problem persists.</p>
            </div>
        </div>
    )
}