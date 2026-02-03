import { useState, useEffect } from 'react';
import { resumeService, authService } from '../services/api';
import { useNavigate } from 'react-router-dom';

export default function Dashboard() {
  const [resumes, setResumes] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [targetRole, setTargetRole] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    loadResumes();
  }, []);

  const loadResumes = async () => {
    try {
      const data = await resumeService.list();
      // Map backend field names (uppercase) to frontend (lowercase)
      const normalizedData = data?.map((resume: any) => ({
        id: resume.ID,
        user_id: resume.UserID,
        title: resume.Title,
        target_role: resume.TargetRole,
        summary: resume.Summary,
        projects: resume.Projects,
        skills: resume.Skills,
        is_default: resume.IsDefault,
        created_at: resume.CreatedAt,
        updated_at: resume.UpdatedAt,
      }));
      setResumes(normalizedData || []);
    } catch (error) {
      console.error('Failed to load resumes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async (e: React.FormEvent) => {
    e.preventDefault();
    setGenerating(true);
    try {
      await resumeService.generate(targetRole);
      setTargetRole('');
      await loadResumes();
    } catch (error) {
      console.error('Failed to generate resume:', error);
      alert('Failed to generate resume');
    } finally {
      setGenerating(false);
    }
  };

  const handleLogout = () => {
    authService.logout();
    navigate('/');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-indigo-600">Resume Builder</h1>
          <button
            onClick={handleLogout}
            className="text-gray-600 hover:text-gray-900"
          >
            Logout
          </button>
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-lg shadow-md p-6 mb-8">
            <h2 className="text-2xl font-semibold mb-4">Generate New Resume</h2>
            <form onSubmit={handleGenerate} className="flex gap-4">
              <input
                type="text"
                value={targetRole}
                onChange={(e) => setTargetRole(e.target.value)}
                placeholder="Target role (e.g., Backend Engineer)"
                className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              />
              <button
                type="submit"
                disabled={generating}
                className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg disabled:opacity-50"
              >
                {generating ? 'Generating...' : 'Generate'}
              </button>
            </form>
          </div>

          <div className="space-y-4">
            <h2 className="text-2xl font-semibold">Your Resumes</h2>
            {resumes.length === 0 ? (
              <div className="bg-white rounded-lg shadow-md p-8 text-center text-gray-500">
                No resumes yet. Generate your first one above!
              </div>
            ) : (
              resumes.map((resume) => (
                <div
                  key={resume.id}
                  className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition cursor-pointer"
                  onClick={() => navigate(`/resume/${resume.id}`)}
                >
                  <h3 className="text-xl font-semibold mb-2">{resume.title || 'Untitled Resume'}</h3>
                  <p className="text-gray-600 mb-2">{resume.target_role || 'No role specified'}</p>
                  <p className="text-sm text-gray-500">
                    Created: {resume.created_at ? new Date(resume.created_at).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' }) : 'Recently'}
                  </p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
