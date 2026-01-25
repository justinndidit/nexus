# Stage 1: Build
FROM python:3.12-slim AS builder

WORKDIR /app

COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

# Stage 2: Run
FROM python:3.12-slim

WORKDIR /app

# Copy only the installed packages from builder
COPY --from=builder /root/.local /root/.local
COPY . .

# Ensure scripts in .local are on the PATH
ENV PATH=/root/.local/bin:$PATH

# Run without --reload for performance and stability
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]