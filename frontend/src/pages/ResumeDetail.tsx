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
      // Normalize field names and nested objects
      const normalized = data ? {
        id: data.ID,
        title: data.Title,
        target_role: data.TargetRole,
        summary: data.Summary,
        projects: data.Projects?.map((p: any) => ({
          repo_name: p.RepoName,
          description: p.Description,
          url: p.URL,
          stars: p.Stars,
          language: p.Language,
          topics: p.Topics,
          highlights: p.Highlights,
        })),
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

  const handleExport = () => {
    const element = document.createElement('a');
    const file = new Blob([JSON.stringify(resume, null, 2)], { type: 'application/json' });
    element.href = URL.createObjectURL(file);
    element.download = `${resume.title.replace(/\s+/g, '_')}_resume.json`;
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  };

  const handlePrint = () => {
    window.print();
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
          <div className="flex gap-4">
            <button
              onClick={handlePrint}
              className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg"
            >
              Print/PDF
            </button>
            <button
              onClick={handleExport}
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg"
            >
              Export JSON
            </button>
            <button
              onClick={handleDelete}
              className="text-red-600 hover:text-red-700 px-4 py-2"
            >
              Delete
            </button>
          </div>
        </div>
      </nav>

      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-lg p-12 print:shadow-none">
          {/* Header */}
          <div className="mb-8 pb-6 border-b-2 border-gray-200">
            <h1 className="text-4xl font-bold mb-2 text-gray-900">{resume.title}</h1>
            <p className="text-2xl text-indigo-600 font-semibold mb-3">{resume.target_role}</p>
            <p className="text-gray-700 leading-relaxed text-lg">{resume.summary}</p>
          </div>

          {/* Skills */}
          <div className="mb-8">
            <h2 className="text-2xl font-bold mb-4 text-gray-900 uppercase tracking-wide">Technical Skills</h2>
            <div className="flex flex-wrap gap-2">
              {resume.skills?.slice(0, 12).map((skill: string, index: number) => (
                <span
                  key={index}
                  className="bg-indigo-50 text-indigo-700 px-4 py-2 rounded-md text-sm font-medium border border-indigo-200"
                >
                  {skill}
                </span>
              ))}
            </div>
          </div>

          {/* Projects */}
          <div>
            <h2 className="text-2xl font-bold mb-6 text-gray-900 uppercase tracking-wide">Key Projects</h2>
            <div className="space-y-6">
              {resume.projects?.map((project: any, index: number) => (
                <div key={index} className="border-l-4 border-indigo-600 pl-6 py-2">
                  <div className="flex items-start justify-between mb-2">
                    <h3 className="text-xl font-bold text-gray-900">
                      <a
                        href={project.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-indigo-600 hover:text-indigo-800 hover:underline"
                      >
                        {project.repo_name}
                      </a>
                    </h3>
                    <div className="flex items-center gap-3 text-sm text-gray-600">
                      {project.stars > 0 && <span className="font-semibold">⭐ {project.stars}</span>}
                      {project.language && (
                        <span className="bg-gray-100 px-2 py-1 rounded text-xs font-medium">
                          {project.language}
                        </span>
                      )}
                    </div>
                  </div>
                  
                  {project.description && (
                    <p className="text-gray-700 mb-3 leading-relaxed">{project.description}</p>
                  )}
                  
                  {project.topics?.length > 0 && (
                    <div className="flex flex-wrap gap-2 mb-2">
                      {project.topics.slice(0, 5).map((topic: string, i: number) => (
                        <span key={i} className="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded">
                          #{topic}
                        </span>
                      ))}
                    </div>
                  )}
                  
                  {project.highlights?.length > 0 && (
                    <ul className="mt-3 space-y-1">
                      {project.highlights.map((highlight: string, i: number) => (
                        <li key={i} className="text-sm text-gray-600 flex items-start">
                          <span className="mr-2 text-indigo-600">▸</span>
                          <span>{highlight}</span>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* Footer */}
          <div className="mt-8 pt-6 border-t border-gray-200 text-center text-sm text-gray-500">
            <p>Generated from GitHub profile • {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
