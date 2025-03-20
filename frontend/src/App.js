import React, { useEffect, useState } from "react";
import axios from "axios";
import "./styles.css";

const App = () => {
  const [tasks, setTasks] = useState([]);
  const [developers, setDevelopers] = useState([]);
  const [assignments, setAssignments] = useState(null);

  useEffect(() => {
    axios.get("http://localhost:8080/tasks")
      .then(response => setTasks(response.data.tasks))
      .catch(error => console.error("Error fetching tasks:", error));

    axios.get("http://localhost:8080/developers")
      .then(response => setDevelopers(response.data.developers))
      .catch(error => console.error("Error fetching developers:", error));
  }, []);

  const handleSchedule = () => {
    axios.get("http://localhost:8080/tasks/schedule")
      .then(response => setAssignments(response.data))
      .catch(error => console.error("Error scheduling tasks:", error));
  };

  return (
    <div className="container">
      <h1>Task Planner</h1>

      <h2>Developers</h2>
      <div className="developer-list">
        {developers.map(dev => (
          <div key={dev.id} className="developer-card">
            <h3>{dev.firstName} {dev.lastName}</h3>
            <p>Email: {dev.email}</p>
            <p>Capacity: {dev.capacity} hrs</p>
          </div>
        ))}
      </div>

      <br />
      <hr />
      <br />

      <h2>Tasks</h2>
      <div className="task-list">
        {tasks.map(task => (
          <div key={task.id} className="task-card">
            <h4>{task.name}</h4>
            <p>Difficulty: {task.difficulty}</p>
            <p>Duration: {task.duration} min</p>
          </div>
        ))}
      </div>

        <br />
      <hr />
      <br />

      <button onClick={handleSchedule} className="schedule-btn">Schedule Tasks</button>

      {assignments && (
        <AssignmentsSummary assignments={assignments} />
      )}
    </div>
  );
};

const AssignmentsSummary = ({ assignments }) => {
  return (
    <div className="assignments-container">
      <h2>Scheduled Tasks</h2>
      <p>Total Work Days: {assignments.totalWorkDay}</p>
      <p>Total Elapsed Work Hours: {assignments.totalElapsedWorkHour}</p>
      <p>Minimum Weeks Required: {assignments.minWeek}</p>

      {assignments.assignments.map((assignment, index) => (
        <div key={index} className="week-container">
          <h3>Week {index + 1}</h3>
          <DeveloperTasks developerTasks={assignment.developerTasks} />
        </div>
      ))}
    </div>
  );
};

const DeveloperTasks = ({ developerTasks }) => {
  return (
    <div className="developer-tasks">
      {developerTasks.map((devTask, index) => (
        <div key={index} className="developer-task-card">
          <h3>{devTask.developer.firstName} {devTask.developer.lastName}</h3>
          <p><strong>Email:</strong> {devTask.developer.email}</p>
          <p><strong>Capacity:</strong> {devTask.developer.capacity}</p>
          <p><strong>Assigned Tasks:</strong></p>
          <ul>
            {devTask.tasks.map(task => (
              <li key={task.id}>
                <strong>{task.name}</strong> - {task.duration}h (Difficulty: {task.difficulty})<br />
                <small>Provider: {task.provider || "N/A"}</small>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
};

export default App;