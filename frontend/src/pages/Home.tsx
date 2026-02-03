import { authService } from '../services/api';

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-5xl font-bold text-gray-900 mb-6">
            GitHub Resume Builder
          </h1>
          <p className="text-xl text-gray-600 mb-8">
            Transform your GitHub profile into a professional resume in seconds
          </p>
          
          <div className="bg-white rounded-lg shadow-xl p-8 mb-12">
            <h2 className="text-2xl font-semibold mb-4">How it works</h2>
            <div className="grid md:grid-cols-3 gap-6 text-left">
              <div className="p-4">
                <div className="text-3xl mb-2">🔐</div>
                <h3 className="font-semibold mb-2">1. Connect GitHub</h3>
                <p className="text-gray-600">Sign in with your GitHub account</p>
              </div>
              <div className="p-4">
                <div className="text-3xl mb-2">🤖</div>
                <h3 className="font-semibold mb-2">2. AI Generation</h3>
                <p className="text-gray-600">Our AI analyzes your repositories</p>
              </div>
              <div className="p-4">
                <div className="text-3xl mb-2">📄</div>
                <h3 className="font-semibold mb-2">3. Get Resume</h3>
                <p className="text-gray-600">Download your professional resume</p>
              </div>
            </div>
          </div>

          <button
            onClick={() => authService.login()}
            className="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-4 px-8 rounded-lg text-lg transition duration-200 shadow-lg"
          >
            Sign in with GitHub
          </button>
        </div>
      </div>
    </div>
  );
}
