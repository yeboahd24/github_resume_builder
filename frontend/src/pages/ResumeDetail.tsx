import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { resumeService } from '../services/api';

export default function ResumeDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [resume, setResume] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadResume();
  }, [id]);

  const loadResume = async () => {
    try {
      const data = await resumeService.get(Number(id));
      // Normalize field names
      const normalized = data ? {
        id: data.ID,
        title: data.Title,
        target_role: data.TargetRole,
        summary: data.Summary,
        projects: data.Projects,
        skills: data.Skills,
      } : null;
      setResume(normalized);
    } catch (error) {
      console.error('Failed to load resume:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this resume?')) {
      try {
        await resumeService.delete(Number(id));
        navigate('/dashboard');
      } catch (error) {
        console.error('Failed to delete resume:', error);
      }
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  if (!resume) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Resume not found</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <button
            onClick={() => navigate('/dashboard')}
            className="text-indigo-600 hover:text-indigo-700"
          >
            ← Back to Dashboard
          </button>
          <button
            onClick={handleDelete}
            className="text-red-600 hover:text-red-700"
          >
            Delete
          </button>
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-lg p-8">
          <div className="mb-8">
            <h1 className="text-3xl font-bold mb-2">{resume.title}</h1>
            <p className="text-xl text-gray-600 mb-4">{resume.target_role}</p>
            <p className="text-gray-700 leading-relaxed">{resume.summary}</p>
          </div>

          <div className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">Skills</h2>
            <div className="flex flex-wrap gap-2">
              {resume.skills?.map((skill: string, index: number) => (
                <span
                  key={index}
                  className="bg-indigo-100 text-indigo-800 px-3 py-1 rounded-full text-sm"
                >
                  {skill}
                </span>
              ))}
            </div>
          </div>

          <div>
            <h2 className="text-2xl font-semibold mb-4">Projects</h2>
            <div className="space-y-6">
              {resume.projects?.map((project: any, index: number) => (
                <div key={index} className="border-l-4 border-indigo-500 pl-4">
                  <h3 className="text-xl font-semibold mb-2">
                    <a
                      href={project.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-indigo-600 hover:underline"
                    >
                      {project.repo_name}
                    </a>
                  </h3>
                  <p className="text-gray-700 mb-2">{project.description}</p>
                  <div className="flex items-center gap-4 text-sm text-gray-600">
                    <span>⭐ {project.stars}</span>
                    <span>{project.language}</span>
                  </div>
                  {project.highlights?.length > 0 && (
                    <ul className="mt-2 space-y-1">
                      {project.highlights.map((highlight: string, i: number) => (
                        <li key={i} className="text-sm text-gray-600">
                          • {highlight}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
