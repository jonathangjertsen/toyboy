class WebGLFrameRenderer {
    constructor(canvas) {
        this.canvas = canvas;
        this.lastCanvasWidth = this.canvas.width;
        this.lastCanvasHeight = this.canvas.height;

        this.gl = canvas.getContext('webgl2') || canvas.getContext('webgl');
        if (!this.gl) {
            throw new Error('WebGL not supported');
        } 
        this.setupWebGL();
    }
   
    setupWebGL() {
        const gl = this.gl;

        const program = gl.createProgram();
        gl.attachShader(program, this.createShader(gl.VERTEX_SHADER, `
            attribute vec2 a_position;
            attribute vec2 a_texCoord;
            varying vec2 v_texCoord;
            void main() {
                gl_Position = vec4(a_position, 0.0, 1.0);
                v_texCoord = a_texCoord;
            }
        `));
        gl.attachShader(program, this.createShader(gl.FRAGMENT_SHADER, `
            precision mediump float;
            uniform sampler2D u_texture;
            varying vec2 v_texCoord;
            void main() {
                gl_FragColor = texture2D(u_texture, v_texCoord);
            }
        `));
        gl.linkProgram(program);
        
        if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
            throw new Error('Program linking failed: ' + gl.getProgramInfoLog(program));
        }
        
        // Map canvas space (-1:1)x(-1:1) to texture space (0:1)x(0:1)
        // Canvas is covered by two triangles
        const vertexBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ARRAY_BUFFER, vertexBuffer);
        gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([
            -1, -1,  0, 1,
            1, -1,  1, 1,
            -1,  1,  0, 0,

            -1,  1,  0, 0,
            1, -1,  1, 1,
            1,  1,  1, 0
        ]), gl.STATIC_DRAW);
        
        const texture = gl.createTexture();
        gl.bindTexture(gl.TEXTURE_2D, texture);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);

        gl.useProgram(program);
        gl.bindTexture(gl.TEXTURE_2D, texture);

        gl.bindBuffer(gl.ARRAY_BUFFER, vertexBuffer);
        
        const aPosition = gl.getAttribLocation(program, 'a_position');
        gl.enableVertexAttribArray(aPosition);
        gl.vertexAttribPointer(aPosition, 2, gl.FLOAT, false, 16, 0);
        
        const aTexCoord = gl.getAttribLocation(program, 'a_texCoord');
        gl.enableVertexAttribArray(aTexCoord);
        gl.vertexAttribPointer(aTexCoord, 2, gl.FLOAT, false, 16, 8);

        const uTexture = gl.getUniformLocation(program, 'u_texture');
        gl.uniform1i(uTexture, 0);
    }

    createShader(type, source) {
        const gl = this.gl;
        const shader = gl.createShader(type);
        gl.shaderSource(shader, source);
        gl.compileShader(shader);
        if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
            const error = gl.getShaderInfoLog(shader);
            gl.deleteShader(shader);
            throw new Error('Shader compilation failed: ' + error);
        }
        return shader;
    }

    
    renderFrame(frameData, width, height) {
        const gl = this.gl;

        // Resize viewport if needed
        if ((this.lastCanvasWidth !== this.canvas.width) || (this.lastCanvasHeight !== this.canvas.height)) {
            gl.viewport(0, 0, this.canvas.width, this.canvas.height);
            this.lastCanvasWidth = this.canvas.width;
            this.lastCanvasHeight = this.canvas.height;
        }

        gl.viewport(0, 0, this.canvas.width, this.canvas.height);
        gl.texImage2D(gl.TEXTURE_2D, 0, gl.LUMINANCE, width, height, 0, gl.LUMINANCE, gl.UNSIGNED_BYTE, frameData);
        gl.drawArrays(gl.TRIANGLES, 0, 6);
    }
}
