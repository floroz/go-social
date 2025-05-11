import apiClient from "@/lib/api";
import type {
  ListPostsSuccessResponse,
  CreatePostRequest,
  CreatePostSuccessResponse,
  // Post, // Not directly used in this file's function signatures yet, but available
} from "@/types/api";

// Local type aliases are no longer needed as types are imported directly.

const PostService = {
  /**
   * Fetches a list of posts from the backend.
   * @returns A promise that resolves with the list of posts.
   */
  listPosts: async (): Promise<ListPostsSuccessResponse> => {
    // The actual response from apiClient.get will be AxiosResponse<ListPostsSuccessResponse>
    // We are interested in response.data which is ListPostsSuccessResponse
    const response = await apiClient.get<ListPostsSuccessResponse>("/v1/posts");
    return response.data; // Contains { data: Post[] }
  },

  /**
   * Creates a new post.
   * @param postData - The data for the new post.
   * @returns A promise that resolves with the created post.
   */
  createPost: async (
    postData: CreatePostRequest
  ): Promise<CreatePostSuccessResponse> => {
    // apiClient.post will take { data: CreatePostRequest } due to convention
    const response = await apiClient.post<CreatePostSuccessResponse>(
      "/v1/posts",
      { data: postData } // Ensure payload is wrapped
    );
    return response.data; // Contains { data: Post }
  },

  // TODO: Add other post-related service methods as needed:
  // getPostById: async (id: number): Promise<GetPostSuccessResponse> => { ... }
  // updatePost: async (id: number, data: UpdatePostRequest): Promise<UpdatePostSuccessResponse> => { ... }
  // deletePost: async (id: number): Promise<void> => { ... }
};

export default PostService;
