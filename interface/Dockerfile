# Use an official Node.js runtime as a parent image
FROM node:18

# Set the working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install any needed packages
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React app
RUN npm run build

# Install serve to serve the build files
RUN npm install -g serve

# Expose the port the app runs on
EXPOSE 5000

# Define the command to run the app
CMD ["serve", "-s", "build", "-l", "5000"]
