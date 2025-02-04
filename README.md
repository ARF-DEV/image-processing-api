# Image Processing API
image processing API, solution for roadmap.sh [project](https://roadmap.sh/projects/image-processing-service)

## Feature
### User Authentication:
- Sign-Up: Allow users to create an account.
- Log-In: Allow users to log into their account.
- JWT Authentication: Secure endpoints using JWTs for authenticated access.

### Image Management
- Upload Image: Allow users to upload images.
- Transform Image: Allow users to perform various transformations (resize, crop, rotate, watermark etc.).
- Retrieve Image: Allow users to retrieve a saved image in different formats.
- List Images: List all uploaded images by the user with metadata.
- Image Transformation

#### List of implemented transformations:
- Resize
- Crop
- Rotate
- Watermark
- Flip
- Mirror
- Compress
- Change format (JPEG and PNG)
- Apply filters (grayscale and sepia)


Note: All expect Rotate transformation, are implemented without any 3rd party libraries (expect for golang's standard library)

## API Specification
1. Register a new user:
```
POST /register
// Request
{
  "username": "user1",
  "password": "password123"
}

// Response
{
	"message": "success",
	"code": "success",
	"data": null,
	"errors": []
}

```


2. Log in an existing user:
```
POST /login
// Request
{
  "username": "user1",
  "password": "password123"
}

// Response
{
	"message": "success",
	"code": "success",
	"data": null
	"errors": []
}
```


3. Upload an image:
```
POST /images
Request Body: Multipart form-data with image file
{
  "image": image_file
}

// Response
{
	"message": "success",
	"code": "success",
	"data": null,
	"errors": []
}
```

4. Apply transformations to an image:
```
POST /images/:id/transform
// Request
{
  "transformations": {
    "resize": {
      "width": "number",
      "height": "number"
    },
    "crop": {
      "width": "number",
      "height": "number",
      "x": "number",
      "y": "number"
    },
    "rotate": "number",
    "format": "string",
    "filters": {
      "grayscale": "boolean",
      "sepia": "boolean"
    }
  }
}
```

5. Retrieve an image:
```
GET /images/:id
// Response
{
	"message": "success",
	"code": "success",
	"data": {
		"id": 1,
		"url": "https://storage.googleapis.com/xxxxxx/749574.jpg"
	},
	"errors": []
}
```

6. Get a paginated list of images:
```
GET /images?page=1&limit=10
// Response
{
	"message": "success",
	"code": "success",
	"data": [
		{
			"id": 1,
			"url": "https://storage.googleapis.com/xxxx/749574.jpg"
		},
		{
			"id": 5,
			"url": "https://storage.googleapis.com/xxxx/749574:resized-rotated-formated-filtered.jpg"
		},
		{
			"id": 6,
			"url": "https://storage.googleapis.com/xxxx/749574:resized-rotated-formated-filtered.jpg"
		},
		{
			"id": 7,
			"url": "https://storage.googleapis.com/xxxx/749574:resized-rotated-formated-filtered.jpg"
		},
		{
			"id": 8,
			"url": "https://storage.googleapis.com/xxxx/749574:resized-rotated-formated-filtered.jpg"
		}
	],
	"meta": {
		"page": 1,
		"limit": 5,
		"total_data": 5,
		"total_page": 1
	},
	"errors": []
}
```
