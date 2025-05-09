PGDMP  (    9        
        }         
   academixDB    17.4    17.4 4    ,           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                           false            -           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                           false            .           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                           false            /           1262    16564 
   academixDB    DATABASE     r   CREATE DATABASE "academixDB" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en-US';
    DROP DATABASE "academixDB";
                     postgres    false            �            1259    16641    assignment_submissions    TABLE     Q  CREATE TABLE public.assignment_submissions (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    assignment_id bigint NOT NULL,
    student_id bigint NOT NULL,
    submission character varying(255),
    marks bigint,
    feedback text
);
 *   DROP TABLE public.assignment_submissions;
       public         heap r       postgres    false            �            1259    16640    assignment_submissions_id_seq    SEQUENCE     �   CREATE SEQUENCE public.assignment_submissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 4   DROP SEQUENCE public.assignment_submissions_id_seq;
       public               postgres    false    226            0           0    0    assignment_submissions_id_seq    SEQUENCE OWNED BY     _   ALTER SEQUENCE public.assignment_submissions_id_seq OWNED BY public.assignment_submissions.id;
          public               postgres    false    225            �            1259    16626    assignments    TABLE     +  CREATE TABLE public.assignments (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    serial bigint,
    course_code text NOT NULL,
    instructions text,
    publish_time timestamp with time zone
);
    DROP TABLE public.assignments;
       public         heap r       postgres    false            �            1259    16625    assignments_id_seq    SEQUENCE     {   CREATE SEQUENCE public.assignments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 )   DROP SEQUENCE public.assignments_id_seq;
       public               postgres    false    224            1           0    0    assignments_id_seq    SEQUENCE OWNED BY     I   ALTER SEQUENCE public.assignments_id_seq OWNED BY public.assignments.id;
          public               postgres    false    223            �            1259    16580    course_models    TABLE        CREATE TABLE public.course_models (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    code text NOT NULL,
    title text NOT NULL,
    description text
);
 !   DROP TABLE public.course_models;
       public         heap r       postgres    false            �            1259    16579    course_models_id_seq    SEQUENCE     }   CREATE SEQUENCE public.course_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 +   DROP SEQUENCE public.course_models_id_seq;
       public               postgres    false    220            2           0    0    course_models_id_seq    SEQUENCE OWNED BY     M   ALTER SEQUENCE public.course_models_id_seq OWNED BY public.course_models.id;
          public               postgres    false    219            �            1259    16591    instructor_courses    TABLE     �   CREATE TABLE public.instructor_courses (
    course_model_id bigint NOT NULL,
    course_model_code text NOT NULL,
    user_model_id bigint NOT NULL
);
 &   DROP TABLE public.instructor_courses;
       public         heap r       postgres    false            �            1259    16608    user_courses    TABLE     �   CREATE TABLE public.user_courses (
    course_model_id bigint NOT NULL,
    course_model_code text NOT NULL,
    user_model_id bigint NOT NULL
);
     DROP TABLE public.user_courses;
       public         heap r       postgres    false            �            1259    16566    user_models    TABLE     &  CREATE TABLE public.user_models (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    username text,
    email text,
    password text NOT NULL,
    role text NOT NULL
);
    DROP TABLE public.user_models;
       public         heap r       postgres    false            �            1259    16565    user_models_id_seq    SEQUENCE     {   CREATE SEQUENCE public.user_models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 )   DROP SEQUENCE public.user_models_id_seq;
       public               postgres    false    218            3           0    0    user_models_id_seq    SEQUENCE OWNED BY     I   ALTER SEQUENCE public.user_models_id_seq OWNED BY public.user_models.id;
          public               postgres    false    217            q           2604    16644    assignment_submissions id    DEFAULT     �   ALTER TABLE ONLY public.assignment_submissions ALTER COLUMN id SET DEFAULT nextval('public.assignment_submissions_id_seq'::regclass);
 H   ALTER TABLE public.assignment_submissions ALTER COLUMN id DROP DEFAULT;
       public               postgres    false    225    226    226            p           2604    16629    assignments id    DEFAULT     p   ALTER TABLE ONLY public.assignments ALTER COLUMN id SET DEFAULT nextval('public.assignments_id_seq'::regclass);
 =   ALTER TABLE public.assignments ALTER COLUMN id DROP DEFAULT;
       public               postgres    false    224    223    224            o           2604    16583    course_models id    DEFAULT     t   ALTER TABLE ONLY public.course_models ALTER COLUMN id SET DEFAULT nextval('public.course_models_id_seq'::regclass);
 ?   ALTER TABLE public.course_models ALTER COLUMN id DROP DEFAULT;
       public               postgres    false    219    220    220            n           2604    16569    user_models id    DEFAULT     p   ALTER TABLE ONLY public.user_models ALTER COLUMN id SET DEFAULT nextval('public.user_models_id_seq'::regclass);
 =   ALTER TABLE public.user_models ALTER COLUMN id DROP DEFAULT;
       public               postgres    false    217    218    218            )          0    16641    assignment_submissions 
   TABLE DATA           �   COPY public.assignment_submissions (id, created_at, updated_at, deleted_at, assignment_id, student_id, submission, marks, feedback) FROM stdin;
    public               postgres    false    226   �C       '          0    16626    assignments 
   TABLE DATA           ~   COPY public.assignments (id, created_at, updated_at, deleted_at, serial, course_code, instructions, publish_time) FROM stdin;
    public               postgres    false    224   mD       #          0    16580    course_models 
   TABLE DATA           i   COPY public.course_models (id, created_at, updated_at, deleted_at, code, title, description) FROM stdin;
    public               postgres    false    220   �D       $          0    16591    instructor_courses 
   TABLE DATA           _   COPY public.instructor_courses (course_model_id, course_model_code, user_model_id) FROM stdin;
    public               postgres    false    221   _E       %          0    16608    user_courses 
   TABLE DATA           Y   COPY public.user_courses (course_model_id, course_model_code, user_model_id) FROM stdin;
    public               postgres    false    222   �E       !          0    16566    user_models 
   TABLE DATA           t   COPY public.user_models (id, created_at, updated_at, deleted_at, name, username, email, password, role) FROM stdin;
    public               postgres    false    218   �E       4           0    0    assignment_submissions_id_seq    SEQUENCE SET     K   SELECT pg_catalog.setval('public.assignment_submissions_id_seq', 2, true);
          public               postgres    false    225            5           0    0    assignments_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.assignments_id_seq', 1, true);
          public               postgres    false    223            6           0    0    course_models_id_seq    SEQUENCE SET     B   SELECT pg_catalog.setval('public.course_models_id_seq', 1, true);
          public               postgres    false    219            7           0    0    user_models_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.user_models_id_seq', 3, true);
          public               postgres    false    217            �           2606    16648 2   assignment_submissions assignment_submissions_pkey 
   CONSTRAINT     p   ALTER TABLE ONLY public.assignment_submissions
    ADD CONSTRAINT assignment_submissions_pkey PRIMARY KEY (id);
 \   ALTER TABLE ONLY public.assignment_submissions DROP CONSTRAINT assignment_submissions_pkey;
       public                 postgres    false    226            �           2606    16633    assignments assignments_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignments_pkey PRIMARY KEY (id);
 F   ALTER TABLE ONLY public.assignments DROP CONSTRAINT assignments_pkey;
       public                 postgres    false    224            z           2606    16587     course_models course_models_pkey 
   CONSTRAINT     d   ALTER TABLE ONLY public.course_models
    ADD CONSTRAINT course_models_pkey PRIMARY KEY (id, code);
 J   ALTER TABLE ONLY public.course_models DROP CONSTRAINT course_models_pkey;
       public                 postgres    false    220    220                       2606    16597 *   instructor_courses instructor_courses_pkey 
   CONSTRAINT     �   ALTER TABLE ONLY public.instructor_courses
    ADD CONSTRAINT instructor_courses_pkey PRIMARY KEY (course_model_id, course_model_code, user_model_id);
 T   ALTER TABLE ONLY public.instructor_courses DROP CONSTRAINT instructor_courses_pkey;
       public                 postgres    false    221    221    221            }           2606    16589 $   course_models uni_course_models_code 
   CONSTRAINT     _   ALTER TABLE ONLY public.course_models
    ADD CONSTRAINT uni_course_models_code UNIQUE (code);
 N   ALTER TABLE ONLY public.course_models DROP CONSTRAINT uni_course_models_code;
       public                 postgres    false    220            t           2606    16577 !   user_models uni_user_models_email 
   CONSTRAINT     ]   ALTER TABLE ONLY public.user_models
    ADD CONSTRAINT uni_user_models_email UNIQUE (email);
 K   ALTER TABLE ONLY public.user_models DROP CONSTRAINT uni_user_models_email;
       public                 postgres    false    218            v           2606    16575 $   user_models uni_user_models_username 
   CONSTRAINT     c   ALTER TABLE ONLY public.user_models
    ADD CONSTRAINT uni_user_models_username UNIQUE (username);
 N   ALTER TABLE ONLY public.user_models DROP CONSTRAINT uni_user_models_username;
       public                 postgres    false    218            �           2606    16614    user_courses user_courses_pkey 
   CONSTRAINT     �   ALTER TABLE ONLY public.user_courses
    ADD CONSTRAINT user_courses_pkey PRIMARY KEY (course_model_id, course_model_code, user_model_id);
 H   ALTER TABLE ONLY public.user_courses DROP CONSTRAINT user_courses_pkey;
       public                 postgres    false    222    222    222            x           2606    16573    user_models user_models_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.user_models
    ADD CONSTRAINT user_models_pkey PRIMARY KEY (id);
 F   ALTER TABLE ONLY public.user_models DROP CONSTRAINT user_models_pkey;
       public                 postgres    false    218            �           1259    16659 %   idx_assignment_submissions_deleted_at    INDEX     n   CREATE INDEX idx_assignment_submissions_deleted_at ON public.assignment_submissions USING btree (deleted_at);
 9   DROP INDEX public.idx_assignment_submissions_deleted_at;
       public                 postgres    false    226            �           1259    16639    idx_assignments_deleted_at    INDEX     X   CREATE INDEX idx_assignments_deleted_at ON public.assignments USING btree (deleted_at);
 .   DROP INDEX public.idx_assignments_deleted_at;
       public                 postgres    false    224            {           1259    16590    idx_course_models_deleted_at    INDEX     \   CREATE INDEX idx_course_models_deleted_at ON public.course_models USING btree (deleted_at);
 0   DROP INDEX public.idx_course_models_deleted_at;
       public                 postgres    false    220            r           1259    16578    idx_user_models_deleted_at    INDEX     X   CREATE INDEX idx_user_models_deleted_at ON public.user_models USING btree (deleted_at);
 .   DROP INDEX public.idx_user_models_deleted_at;
       public                 postgres    false    218            �           2606    16649 ;   assignment_submissions fk_assignment_submissions_assignment    FK CONSTRAINT     �   ALTER TABLE ONLY public.assignment_submissions
    ADD CONSTRAINT fk_assignment_submissions_assignment FOREIGN KEY (assignment_id) REFERENCES public.assignments(id) ON DELETE CASCADE;
 e   ALTER TABLE ONLY public.assignment_submissions DROP CONSTRAINT fk_assignment_submissions_assignment;
       public               postgres    false    4739    226    224            �           2606    16634 (   assignments fk_course_models_assignments    FK CONSTRAINT     �   ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT fk_course_models_assignments FOREIGN KEY (course_code) REFERENCES public.course_models(code);
 R   ALTER TABLE ONLY public.assignments DROP CONSTRAINT fk_course_models_assignments;
       public               postgres    false    4733    224    220            �           2606    16598 5   instructor_courses fk_instructor_courses_course_model    FK CONSTRAINT     �   ALTER TABLE ONLY public.instructor_courses
    ADD CONSTRAINT fk_instructor_courses_course_model FOREIGN KEY (course_model_id, course_model_code) REFERENCES public.course_models(id, code);
 _   ALTER TABLE ONLY public.instructor_courses DROP CONSTRAINT fk_instructor_courses_course_model;
       public               postgres    false    221    4730    220    221    220            �           2606    16603 3   instructor_courses fk_instructor_courses_user_model    FK CONSTRAINT     �   ALTER TABLE ONLY public.instructor_courses
    ADD CONSTRAINT fk_instructor_courses_user_model FOREIGN KEY (user_model_id) REFERENCES public.user_models(id);
 ]   ALTER TABLE ONLY public.instructor_courses DROP CONSTRAINT fk_instructor_courses_user_model;
       public               postgres    false    221    4728    218            �           2606    16615 )   user_courses fk_user_courses_course_model    FK CONSTRAINT     �   ALTER TABLE ONLY public.user_courses
    ADD CONSTRAINT fk_user_courses_course_model FOREIGN KEY (course_model_id, course_model_code) REFERENCES public.course_models(id, code);
 S   ALTER TABLE ONLY public.user_courses DROP CONSTRAINT fk_user_courses_course_model;
       public               postgres    false    220    222    220    222    4730            �           2606    16620 '   user_courses fk_user_courses_user_model    FK CONSTRAINT     �   ALTER TABLE ONLY public.user_courses
    ADD CONSTRAINT fk_user_courses_user_model FOREIGN KEY (user_model_id) REFERENCES public.user_models(id);
 Q   ALTER TABLE ONLY public.user_courses DROP CONSTRAINT fk_user_courses_user_model;
       public               postgres    false    4728    218    222            �           2606    16654 1   assignment_submissions fk_user_models_submissions    FK CONSTRAINT     �   ALTER TABLE ONLY public.assignment_submissions
    ADD CONSTRAINT fk_user_models_submissions FOREIGN KEY (student_id) REFERENCES public.user_models(id);
 [   ALTER TABLE ONLY public.assignment_submissions DROP CONSTRAINT fk_user_models_submissions;
       public               postgres    false    218    226    4728            )   �   x���=�0���Wt����M������h������o)�p�q�=�����6�g��A]3���IC#P�xL�{js)��e�q.�"Ft�ZmY;5��W?4����!5�o�e���;�
6R��C�V�C�i~�p���#()��>�      '   d   x�3�4202�5 "3c+#3+CK=KKcm3�1~�������F��9��E��Ź
�99
I�
�y���)覘[Z�X���M����� ��      #   n   x�3�4202�5 "3c+#c+C=c3m3������������������9H2Ə�9���Ȑӿ �(�$3/]!���$5���'3��B!8#3;#1/_�#��+F��� l��      $      x�3�tv562�4����� K�      %      x�3�tv562�4����� D�      !   X  x���Ko�@�5�
��a����
RE����CP+}5���3��M�IN>$a��Gx=�@b!fi&��Ԑ� �?�ҕ:��Bed�R��^'i&�[��ĸ�<݁`/��-[s��=B/��w�������ѮԾ1�8뙍wq����|Z?�����3 ��@���C�P@c��=.�0��{8"?�K��!�N�lP����!������@oO����E\v��"�RѬ���f([�jC #�Mv/GP��"��:��A�r?�A�4Y{V9*��
?�T�gY�d�_�]�z�����./�d�����H� �^g�X*�}y	dY~8��(     